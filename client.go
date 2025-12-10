package httpx

import (
	"net"
	"net/http"
	"time"
)

// client is the concrete implementation of the Client interface.
// It wraps Go's http.Client and applies httpx-level configuration such as
// connection pooling, timeouts, and global request headers.
type client struct {
	httpClient *http.Client // underlying HTTP engine
	Config                  // global configuration settings
}

// Config defines optional settings used when constructing a new httpx client.
// All fields are optional; zero values trigger library defaults.
//
// Typical usage:
//
//	client := httpx.New(&httpx.Config{
//	    RequestTimeout:     5 * time.Second,
//	    ConnectionTimeout:  1 * time.Second,
//	    MaxIdleConnections: 10,
//	    Headers: http.Header{
//	        "Authorization": []string{"Bearer token"},
//	    },
//	})
type Config struct {
	// Headers applied to every request unless overridden by per-request options.
	Headers http.Header

	// MaxIdleConnections controls the number of idle TCP connections kept per host.
	// Increasing this improves throughput for high-volume APIs.
	MaxIdleConnections int

	// ConnectionTimeout defines the timeout for establishing new TCP connections.
	// A value of 0 disables the timeout and uses default system behavior.
	ConnectionTimeout time.Duration

	// RequestTimeout sets a maximum duration for the full request, including
	// connecting, redirects, and reading the response body. A value of 0
	// disables the timeout entirely.
	RequestTimeout time.Duration
}

// New constructs and returns a new httpx client.
// Missing or zero-valued configuration fields are replaced by defaults.
//
// When cfg is nil, all defaults are applied.
func New(cfg *Config) Client {

	// Apply default settings
	defaults := &Config{
		MaxIdleConnections: 5,
		ConnectionTimeout:  0,
		RequestTimeout:     0,
		Headers:            make(http.Header),
	}

	// Override defaults with user config
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

	// Build the underlying http.Client
	httpClient := &http.Client{
		Timeout: defaults.RequestTimeout, // total request timeout

		Transport: &http.Transport{
			MaxIdleConnsPerHost:   defaults.MaxIdleConnections,
			ResponseHeaderTimeout: defaults.RequestTimeout,

			// TCP dialer configuration
			DialContext: (&net.Dialer{
				Timeout: defaults.ConnectionTimeout,
			}).DialContext,
		},
	}

	return &client{
		httpClient: httpClient,
		Config:     *defaults,
	}
}

// Get performs an HTTP GET request.
// Optional per-request settings can be applied using Option functions.
//
// Example:
//
//	res, err := client.Get("https://api.com/items",
//	    httpx.WithParams(map[string]string{"limit": "10"}))
func (c *client) Get(url string, opts ...Option) (*http.Response, error) {
	return c.do(http.MethodGet, url, buildOptions(opts))
}

// Post performs an HTTP POST request.
// Body encoding is based on Content-Type (JSON, XML, form, multipart, etc.)
//
// Example:
//
//	res, err := client.Post("https://api.com/users",
//	    httpx.WithJSON(user))
func (c *client) Post(url string, opts ...Option) (*http.Response, error) {
	return c.do(http.MethodPost, url, buildOptions(opts))
}

// Put performs an HTTP PUT request.
// Typically used for complete resource replacement.
func (c *client) Put(url string, opts ...Option) (*http.Response, error) {
	return c.do(http.MethodPut, url, buildOptions(opts))
}

// Patch performs an HTTP PATCH request.
// Typically used for partial resource updates.
func (c *client) Patch(url string, opts ...Option) (*http.Response, error) {
	return c.do(http.MethodPatch, url, buildOptions(opts))
}

// Delete performs an HTTP DELETE request.
// DELETE bodies are intentionally not supported to avoid inconsistent behavior
// across HTTP servers.
func (c *client) Delete(url string, opts ...Option) (*http.Response, error) {
	return c.do(http.MethodDelete, url, buildOptions(opts))
}
