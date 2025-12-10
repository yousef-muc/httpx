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
//
// It applies global headers, merges per-request overrides, encodes the request
// body based on Content-Type, appends query parameters, and finally executes the
// HTTP request using the underlying *http.Client.
//
// This method is not exposed publicly; the public API consists of Get, Post,
// Put, Patch, and Delete.
func (c *client) do(method, uri string, o *RequestOptions) (*http.Response, error) {

	//────────────────────────────────────────────────────────────
	// Merge global headers with per-request headers
	//────────────────────────────────────────────────────────────
	requestHeaders := make(http.Header)

	// Apply global headers (from Config)
	for key, values := range c.Headers {
		if len(values) > 0 {
			requestHeaders.Set(key, values[0])
		}
	}

	// Override with per-request headers (from options)
	if o.Headers != nil {
		for key, values := range o.Headers {
			if len(values) > 0 {
				requestHeaders.Set(key, values[0])
			}
		}
	}

	//────────────────────────────────────────────────────────────
	// Validate body usage
	//────────────────────────────────────────────────────────────
	if method == http.MethodGet && o.Body != nil {
		return nil, fmt.Errorf("GET request cannot contain a body")
	}

	// Assign default Content-Type if a body exists but user didn't specify one.
	if o.Body != nil && requestHeaders.Get("Content-Type") == "" {
		requestHeaders.Set("Content-Type", "application/json")
	}

	// Determine base Content-Type (strip charset or options)
	contentType := strings.ToLower(strings.Split(requestHeaders.Get("Content-Type"), ";")[0])

	//────────────────────────────────────────────────────────────
	// Encode request body
	//────────────────────────────────────────────────────────────
	var requestBody []byte

	if o.Body != nil && method != http.MethodGet {
		var err error

		switch contentType {

		// JSON ----------------------------------------------------
		case "application/json":
			requestBody, err = json.Marshal(o.Body)

		// FORM URLENCODED -----------------------------------------
		case "application/x-www-form-urlencoded":
			values := url.Values{}

			switch v := o.Body.(type) {
			case map[string]string:
				for k, val := range v {
					values.Set(k, val)
				}
			case url.Values:
				values = v
			default:
				return nil, fmt.Errorf("body must be map[string]string or url.Values for x-www-form-urlencoded")
			}

			requestBody = []byte(values.Encode())

		// XML -----------------------------------------------------
		case "application/xml", "text/xml":
			requestBody, err = xml.Marshal(o.Body)

		// MULTIPART FORM DATA -------------------------------------
		case "multipart/form-data":
			var b bytes.Buffer
			writer := multipart.NewWriter(&b)

			// Automatically set boundary in Content-Type
			requestHeaders.Set("Content-Type", writer.FormDataContentType())

			fields, ok := o.Body.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("multipart/form-data requires body = map[string]any")
			}

			for key, val := range fields {
				switch cast := val.(type) {

				case []byte:
					// file upload (raw bytes)
					part, err := writer.CreateFormFile(key, key)
					if err != nil {
						return nil, err
					}
					if _, err := part.Write(cast); err != nil {
						return nil, err
					}

				case string:
					// form field value
					if err := writer.WriteField(key, cast); err != nil {
						return nil, err
					}

				default:
					return nil, fmt.Errorf("unsupported multipart field type %T for key %s", cast, key)
				}
			}

			writer.Close()
			requestBody = b.Bytes()

		// PLAIN TEXT ----------------------------------------------
		case "text/plain":
			requestBody = []byte(fmt.Sprintf("%v", o.Body))

		// RAW STREAM / BYTES --------------------------------------
		case "application/octet-stream":
			switch v := o.Body.(type) {
			case []byte:
				requestBody = v
			case io.Reader:
				requestBody, err = io.ReadAll(v)
			default:
				return nil, fmt.Errorf("octet-stream requires []byte or io.Reader body")
			}

		// DEFAULT → JSON ------------------------------------------
		default:
			requestBody, err = json.Marshal(o.Body)
		}

		if err != nil {
			return nil, err
		}
	}

	//────────────────────────────────────────────────────────────
	// Wrap encoded body in an io.Reader
	//────────────────────────────────────────────────────────────
	var bodyReader io.Reader
	if requestBody != nil {
		bodyReader = bytes.NewBuffer(requestBody)
	}

	//────────────────────────────────────────────────────────────
	// Append query parameters (?key=value)
	//────────────────────────────────────────────────────────────
	if o.Params != nil {
		u, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}

		q := u.Query()
		for key, val := range o.Params {
			q.Set(key, val)
		}

		u.RawQuery = q.Encode()
		uri = u.String()
	}

	//────────────────────────────────────────────────────────────
	// Construct the *http.Request
	//────────────────────────────────────────────────────────────
	req, err := http.NewRequest(method, uri, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header = requestHeaders

	//────────────────────────────────────────────────────────────
	// Execute request using the underlying http.Client
	//────────────────────────────────────────────────────────────
	return c.httpClient.Do(req)
}
