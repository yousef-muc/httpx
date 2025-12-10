package httpx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// HttpError represents an HTTP error returned by the server when the response
// status code is outside the 2xx range. It includes contextual information
// such as the status code, response body, headers, and the originating request.
// This mirrors the behavior of Axios' error.response object.
type HttpError struct {
	StatusCode int         // HTTP status code (e.g., 404, 500)
	Status     string      // Full status string (e.g., "404 Not Found")
	Body       []byte      // Raw response body for debugging or custom decoding
	Headers    http.Header // Response headers returned by the server
	Method     string      // HTTP method of the originating request
	URL        string      // Request URL that caused the error
}

// Error implements the error interface. A short body snippet is included
// to make debugging easier without overwhelming logs.
func (e *HttpError) Error() string {
	snippet := string(e.Body)
	if len(snippet) > 200 {
		snippet = snippet[:200] + "..."
	}
	return fmt.Sprintf("httpx: %s %s returned %d (%s)", e.Method, e.URL, e.StatusCode, snippet)
}

// readBodyWithStatus reads and returns the full response body. If the response
// status code is not within the 2xx success range, an HttpError is returned
// containing the response metadata.
// This function is used internally by all response helpers.
func readBodyWithStatus(res *http.Response) ([]byte, error) {
	defer res.Body.Close()

	// Read raw body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Non-2xx responses return an HttpError
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, &HttpError{
			StatusCode: res.StatusCode,
			Status:     res.Status,
			Body:       body,
			Headers:    res.Header.Clone(),
			Method:     res.Request.Method,
			URL:        res.Request.URL.String(),
		}
	}

	return body, nil
}

// Bytes reads and returns the response body as raw bytes. If the response
// contains a non-2xx status code, an HttpError is returned instead.
func (c *client) Bytes(res *http.Response) ([]byte, error) {
	return readBodyWithStatus(res)
}

// Text reads and returns the response body as a UTF-8 string.
// Non-2xx responses return an HttpError.
func (c *client) Text(res *http.Response) (string, error) {
	b, err := readBodyWithStatus(res)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadJSON decodes the response body into the provided target struct.
// It returns an HttpError for non-2xx responses or a decoding error
// if the JSON payload cannot be unmarshaled into the target.
//
// Example:
//
//	var user User
//	err := client.ReadJSON(res, &user)
func (c *client) ReadJSON(res *http.Response, target any) error {
	b, err := readBodyWithStatus(res)
	if err != nil {
		return err
	}

	if len(b) == 0 {
		return nil
	}

	if err := json.Unmarshal(b, target); err != nil {
		return fmt.Errorf("httpx: failed to decode JSON: %w", err)
	}

	return nil
}

// JSON decodes a JSON response body into a generic Go type T.
// It returns an error if the status code is non-2xx or if decoding fails.
//
// Example:
//
//	user, err := httpx.JSON[User](res)
func JSON[T any](res *http.Response) (T, error) {
	var out T

	b, err := readBodyWithStatus(res)
	if err != nil {
		return out, err
	}

	if len(b) == 0 {
		return out, nil
	}

	if err := json.Unmarshal(b, &out); err != nil {
		return out, fmt.Errorf("httpx: failed to decode JSON: %w", err)
	}

	return out, nil
}

// ReadXML decodes an XML response body into the provided target struct.
// An HttpError is returned if the status code is not successful.
//
// Example:
//
//	var feed AtomFeed
//	err := client.ReadXML(res, &feed)
func (c *client) ReadXML(res *http.Response, target any) error {
	b, err := readBodyWithStatus(res)
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(b, target); err != nil {
		return fmt.Errorf("httpx: failed to decode XML: %w", err)
	}

	return nil
}

// XML decodes an XML response body into a generic Go type T.
// This is the XML counterpart to JSON[T].
//
// Example:
//
//	feed, err := httpx.XML[Feed](res)
func XML[T any](res *http.Response) (T, error) {
	var out T

	b, err := readBodyWithStatus(res)
	if err != nil {
		return out, err
	}

	if err := xml.Unmarshal(b, &out); err != nil {
		return out, fmt.Errorf("httpx: failed to decode XML: %w", err)
	}

	return out, nil
}
