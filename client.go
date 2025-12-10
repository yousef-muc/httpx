package httpx

import (
	"net"
	"net/http"
	"time"
)

// client is the default implementation of the Client interface.
// It wraps Go's http.Client and holds the configuration used to
// construct transport settings and request defaults.
type client struct {
	httpClient *http.Client
	Config
}

// Config defines optional settings for customizing the underlying
// HTTP client's behavior, including timeouts, connection pooling,
// and default request headers.
type Config struct {
	// Headers are applied to every request unless overridden by per-request headers.
	Headers http.Header

	// MaxIdleConnections controls how many idle TCP connections are kept per host.
	// Higher values improve throughput in high-concurrency applications.
	MaxIdleConnections int

	// ConnectionTimeout defines how long to wait when establishing a TCP connection.
	// A value of 0 disables the timeout and uses Go's default behavior.
	ConnectionTimeout time.Duration

	// RequestTimeout defines the maximum allowed duration for the entire request,
	// including connection establishment, redirects, and reading the response body.
	// A value of 0 disables the timeout.
	RequestTimeout time.Duration
}

// New creates and returns a new httpx client using the provided Config.
// Missing settings are populated using library defaults.
//
//	If cfg is nil, all default values will be used.
func New(cfg *Config) Client {
	// Load default configuration
	defaults := &Config{
		MaxIdleConnections: 5,
		ConnectionTimeout:  0,
		RequestTimeout:     0,
		Headers:            make(http.Header),
	}

	// Override defaults with user-provided configuration
	if cfg != nil {
		if cfg.MaxIdleConnections != 0 {
			defaults.MaxIdleConnections = cfg.MaxIdleConnections
		}
		if cfg.ConnectionTimeout != 0 {
			defaults.ConnectionTimeout = cfg.ConnectionTimeout
		}
		if cfg.RequestTimeout != 0 {
			defaults.RequestTimeout = cfg.RequestTimeout
		}
		if cfg.Headers != nil {
			defaults.Headers = cfg.Headers
		}
	}

	// Construct the underlying http.Client with timeouts and custom transport.
	return &client{
		httpClient: &http.Client{
			// Total request timeout
			Timeout: defaults.RequestTimeout,

			Transport: &http.Transport{
				// Number of idle connections kept per host
				MaxIdleConnsPerHost: defaults.MaxIdleConnections,

				// Timeout for waiting on response headers
				ResponseHeaderTimeout: defaults.RequestTimeout,

				// Dialer with connection timeout
				DialContext: (&net.Dialer{
					Timeout: defaults.ConnectionTimeout,
				}).DialContext,
			},
		},

		// Store final configuration on the client instance
		Config: *defaults,
	}
}

// Get performs an HTTP GET request to the given URL with optional headers
// and query parameters. GET requests cannot contain a request body.
func (c *client) Get(url string, headers http.Header, params map[string]string) (*http.Response, error) {
	return c.do(http.MethodGet, url, headers, params, nil)
}

// Post performs an HTTP POST request using optional headers, query parameters,
// and a request body. The Content-Type determines how the body is encoded.
func (c *client) Post(url string, headers http.Header, params map[string]string, body any) (*http.Response, error) {
	return c.do(http.MethodPost, url, headers, params, body)
}

// Put performs an HTTP PUT request using optional headers, query parameters,
// and a request body. PUT is typically used for full resource replacement.
func (c *client) Put(url string, headers http.Header, params map[string]string, body any) (*http.Response, error) {
	return c.do(http.MethodPut, url, headers, params, body)
}

// Patch performs an HTTP PATCH request using optional headers, query parameters,
// and a request body. PATCH is typically used for partial resource updates.
func (c *client) Patch(url string, headers http.Header, params map[string]string, body any) (*http.Response, error) {
	return c.do(http.MethodPatch, url, headers, params, body)
}

// Delete performs an HTTP DELETE request with optional headers and query parameters.
// DELETE requests may include a body depending on the API, but httpx does not
// support bodies for DELETE calls to avoid inconsistent server behavior.
func (c *client) Delete(url string, headers http.Header, params map[string]string) (*http.Response, error) {
	return c.do(http.MethodDelete, url, headers, params, nil)
}
