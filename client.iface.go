package httpx

import (
	"net/http"
)

// Client defines the public-facing interface for the httpx HTTP client.
//
// It provides a clean, Axios-style API while staying fully compatible with Go's
// native net/http primitives. All request methods support optional request
// configuration through functional options (Option), allowing callers to set
// per-request headers, query parameters, and request bodies.
//
// Each method returns the raw *http.Response object, giving callers full control
// over streaming, manual decoding, or using the response helpers provided by
// httpx (JSON, XML, Text, Bytes).
type Client interface {

	// Get performs an HTTP GET request to the given URL.
	//
	// GET requests cannot include a request body. Optional behavior such as
	// query parameters or additional headers can be configured via Option.
	//
	// Example:
	//    res, err := client.Get("https://api.com/users",
	//        httpx.WithParams(map[string]string{"limit": "10"}),
	//        httpx.WithHeaders(http.Header{"Accept": []string{"application/json"}}),
	//    )
	Get(url string, opts ...Option) (*http.Response, error)

	// Post performs an HTTP POST request using optional headers, query parameters,
	// and a request body. Body encoding is determined automatically based on the
	// Content-Type header (JSON, XML, x-www-form-urlencoded, multipart/form-data, etc.).
	//
	// Example:
	//    res, err := client.Post("https://api.com/users",
	//        httpx.WithBody(user),
	//        httpx.WithHeaders(http.Header{"Content-Type": []string{"application/json"}}),
	//    )
	Post(url string, opts ...Option) (*http.Response, error)

	// Put performs an HTTP PUT request and supports optional headers, parameters,
	// and a body. PUT is generally used for full resource replacement.
	//
	// Example:
	//    res, err := client.Put("https://api.com/users/1",
	//        httpx.WithBody(updatedUser),
	//    )
	Put(url string, opts ...Option) (*http.Response, error)

	// Patch performs an HTTP PATCH request with optional headers, parameters,
	// and a body. PATCH is typically used for partial updates.
	//
	// Example:
	//    res, err := client.Patch("https://api.com/users/1",
	//        httpx.WithBody(map[string]any{"lastname": "Updated"}),
	//    )
	Patch(url string, opts ...Option) (*http.Response, error)

	// Delete performs an HTTP DELETE request. It supports optional headers and
	// query parameters but does not allow request bodies in httpx to avoid
	// inconsistent behavior across HTTP servers.
	//
	// Example:
	//    res, err := client.Delete("https://api.com/users/1",
	//        httpx.WithParams(map[string]string{"force": "true"}),
	//    )
	Delete(url string, opts ...Option) (*http.Response, error)
}
