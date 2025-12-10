package httpx

import (
	"net/http"
)

// Client defines the public interface for performing HTTP requests using httpx.
// It provides a simplified, Axios-like API built on top of Go's net/http package,
// supporting custom headers, query parameters, and automatic request body encoding.
//
// All methods return the raw *http.Response object, allowing callers to decide
// whether to use the provided response helpers (JSON, XML, Text, Bytes) or to
// read and handle the body manually.
type Client interface {

	// Get sends an HTTP GET request to the specified URL.
	//
	// Headers: optional per-request headers; they override global headers.
	// Params:  optional query parameters appended to the URL.
	//
	// GET requests cannot contain a body.
	Get(url string, headers http.Header, params map[string]string) (*http.Response, error)

	// Post sends an HTTP POST request to the specified URL.
	//
	// Headers: optional per-request headers; Content-Type determines encoding.
	// Params:  optional query parameters appended to the URL.
	// Body:    request payload; encoded based on Content-Type (JSON, XML, form, multipart, etc.).
	Post(url string, headers http.Header, params map[string]string, body any) (*http.Response, error)

	// Put sends an HTTP PUT request to the specified URL.
	//
	// PUT typically performs a full resource replacement.
	// Encoding behavior is identical to Post().
	Put(url string, headers http.Header, params map[string]string, body any) (*http.Response, error)

	// Patch sends an HTTP PATCH request to the specified URL.
	//
	// PATCH typically performs a partial resource update.
	// Encoding behavior is identical to Post().
	Patch(url string, headers http.Header, params map[string]string, body any) (*http.Response, error)

	// Delete sends an HTTP DELETE request to the specified URL.
	//
	// Headers: optional per-request headers.
	// Params:  optional query parameters appended to the URL.
	//
	// DELETE requests do not accept a body in httpx for consistency and
	// to avoid differing server behavior across APIs.
	Delete(url string, headers http.Header, params map[string]string) (*http.Response, error)
}
