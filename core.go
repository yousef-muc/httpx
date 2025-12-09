package httpx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"strings"
)

func (c *client) do(method, url string, headers http.Header, body any) (*http.Response, error) {
	client := &http.Client{}

	requestHeaders := c.getRequestHeaders(headers)
	requestBody, err := c.getRequestBody(requestHeaders.Get("Content-Type"), body)
	if err != nil {
		return nil, errors.New("unable to create a request body")
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.New("unable to create a new request")
	}

	request.Header = requestHeaders

	return client.Do(request)
}

func (c *client) getRequestHeaders(reqHeaders http.Header) http.Header {
	result := make(http.Header)
	// add common headers
	for header, values := range c.Headers {
		if len(values) > 0 {
			result.Set(header, values[0])
		}
	}

	// add custom headers
	for header, values := range reqHeaders {
		if len(values) > 0 {
			result.Set(header, values[0])
		}
	}
	return result
}

func (c *client) getRequestBody(contentType string, body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	switch strings.ToLower(contentType) {
	case "application/json":
		return json.Marshal(body)
	case "application/xml":
		return xml.Marshal(body)
	default:
		return json.Marshal(body)
	}
}
