package httpx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// do is the internal request executor used by all HTTP verb methods.
// It assembles headers, query parameters, the request body, and executes the
// request using the underlying http.Client.
//
// Parameters:
//   - method: the HTTP verb (GET, POST, PUT, PATCH, DELETE)
//   - uri: the full request URL (query parameters may be appended)
//   - headers: request-scoped headers that override default client headers
//   - params: optional query parameters
//   - body: the request payload; encoded based on Content-Type
func (c *client) do(method, uri string, headers http.Header, params map[string]string, body any) (*http.Response, error) {

	//────────────────────────────────────────────────────────────
	// Merge global headers with request-specific headers
	//────────────────────────────────────────────────────────────
	requestHeaders := make(http.Header)

	// Apply client-wide default headers
	for header, values := range c.Headers {
		if len(values) > 0 {
			requestHeaders.Set(header, values[0])
		}
	}

	// Override with per-request headers
	for header, values := range headers {
		if len(values) > 0 {
			requestHeaders.Set(header, values[0])
		}
	}

	//────────────────────────────────────────────────────────────
	// Validate body usage
	//────────────────────────────────────────────────────────────

	// GET requests must not contain a body
	if method == http.MethodGet && body != nil {
		return nil, fmt.Errorf("GET request cannot contain a body")
	}

	// Use JSON as the default Content-Type when a body is provided
	if body != nil && requestHeaders.Get("Content-Type") == "" {
		requestHeaders.Set("Content-Type", "application/json")
	}

	// Normalize Content-Type (strip charset, etc.)
	contentType := strings.ToLower(strings.Split(requestHeaders.Get("Content-Type"), ";")[0])

	//────────────────────────────────────────────────────────────
	// Encode request body based on Content-Type
	//────────────────────────────────────────────────────────────
	var requestBody []byte

	if body != nil && method != http.MethodGet {
		var err error

		switch contentType {

		// JSON encoding
		case "application/json":
			requestBody, err = json.Marshal(body)

		// HTML form encoding: expect map[string]string or url.Values
		case "application/x-www-form-urlencoded":
			values := url.Values{}
			switch v := body.(type) {
			case map[string]string:
				for key, val := range v {
					values.Set(key, val)
				}
			case url.Values:
				values = v
			default:
				return nil, fmt.Errorf("body must be map[string]string or url.Values for form-urlencoded requests")
			}
			requestBody = []byte(values.Encode())

		// XML encoding
		case "application/xml", "text/xml":
			requestBody, err = xml.Marshal(body)

		// Multipart form-data (supports text fields + raw file bytes)
		case "multipart/form-data":
			var b bytes.Buffer
			writer := multipart.NewWriter(&b)

			// Set multipart boundary
			requestHeaders.Set("Content-Type", writer.FormDataContentType())

			fields, ok := body.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("body must be map[string]any for multipart/form-data requests")
			}

			for key, val := range fields {
				switch cast := val.(type) {

				case []byte:
					// File upload
					part, err := writer.CreateFormFile(key, key)
					if err != nil {
						return nil, err
					}
					if _, err = part.Write(cast); err != nil {
						return nil, err
					}

				case string:
					// Text field
					if err := writer.WriteField(key, cast); err != nil {
						return nil, err
					}

				default:
					return nil, fmt.Errorf("unsupported multipart field type %T for key %s", val, key)
				}
			}

			writer.Close()
			requestBody = b.Bytes()

		// Plain text encoding
		case "text/plain":
			requestBody = []byte(fmt.Sprintf("%v", body))

		// Raw byte stream: []byte or io.Reader
		case "application/octet-stream":
			switch v := body.(type) {
			case []byte:
				requestBody = v
			case io.Reader:
				requestBody, err = io.ReadAll(v)
			default:
				return nil, fmt.Errorf("octet-stream body must be []byte or io.Reader")
			}

		// Fallback to JSON for unknown Content-Types
		default:
			requestBody, err = json.Marshal(body)
		}

		if err != nil {
			return nil, err
		}
	}

	//────────────────────────────────────────────────────────────
	// Wrap body in an io.Reader
	//────────────────────────────────────────────────────────────
	var bodyReader io.Reader
	if requestBody != nil {
		bodyReader = bytes.NewBuffer(requestBody)
	}

	//────────────────────────────────────────────────────────────
	// Append query parameters to URL
	//────────────────────────────────────────────────────────────
	if params != nil {
		u, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}

		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}

		u.RawQuery = q.Encode()
		uri = u.String()
	}

	//────────────────────────────────────────────────────────────
	// Construct the HTTP request
	//────────────────────────────────────────────────────────────
	request, err := http.NewRequest(method, uri, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	request.Header = requestHeaders

	//────────────────────────────────────────────────────────────
	// Execute request using the underlying http.Client
	//────────────────────────────────────────────────────────────
	return c.httpClient.Do(request)
}
