package httpx

import "net/http"

// RequestOptions holds all optional, per-request configuration.
//
// These values supplement or override global client-level configuration.
// They are populated through functional Option modifiers passed to each
// request method (Get, Post, Put, Patch, Delete).
type RequestOptions struct {
	// Headers are applied only to this specific request and override
	// any global headers defined in the client configuration.
	Headers http.Header

	// Params represent query parameters appended to the request URL.
	// Example: ?page=1&limit=10
	Params map[string]string

	// Body is the request payload. If provided, the Content-Type header
	// determines how the body will be encoded (JSON, XML, form, etc.).
	// GET requests must not include a body.
	Body any
}

// Option is a functional modifier that mutates the RequestOptions struct.
//
// It allows a clean, flexible API such as:
//
//	client.Get(url,
//	    httpx.WithHeaders(h),
//	    httpx.WithParams(p),
//	)
//
// Options can be composed in any order.
type Option func(*RequestOptions)

// WithHeaders applies per-request headers.
//
// These headers override any global headers set in the Config.
// If the same header key exists globally and locally, the local one wins.
//
// Example:
//
//	client.Get(url, httpx.WithHeaders(http.Header{
//	    "Accept": []string{"application/json"},
//	}))
func WithHeaders(h http.Header) Option {
	return func(o *RequestOptions) {
		o.Headers = h
	}
}

// WithParams appends URL query parameters for this request.
//
// Example:
//
//	client.Get(url, httpx.WithParams(map[string]string{
//	    "page": "1",
//	    "limit": "20",
//	}))
func WithParams(p map[string]string) Option {
	return func(o *RequestOptions) {
		o.Params = p
	}
}

// WithBody assigns the request body used by POST, PUT, and PATCH requests.
// GET requests must not include a body and will result in an error.
//
// Example:
//
//	client.Post(url,
//	    httpx.WithBody(myPayload),
//	    httpx.WithHeaders(http.Header{"Content-Type": []string{"application/json"}}),
//	)
func WithBody(b any) Option {
	return func(o *RequestOptions) {
		o.Body = b
	}
}

// buildOptions merges a variadic slice of Option functions into a new
// RequestOptions struct. Missing fields are initialized with sane defaults.
//
// This helper is used internally by all client request methods.
func buildOptions(opts []Option) *RequestOptions {
	o := &RequestOptions{}

	// Apply functional options
	for _, fn := range opts {
		fn(o)
	}

	// Ensure headers map is always initialized
	if o.Headers == nil {
		o.Headers = make(http.Header)
	}

	return o
}
