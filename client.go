package httpx

import "net/http"

type HttpClient interface {
	SetHeaders(headers http.Header)
	Get(string, http.Header) (*http.Response, error)
	Post(url string, headers http.Header, body any) (*http.Response, error)
	Put(url string, headers http.Header, body any) (*http.Response, error)
	Patch(url string, headers http.Header, body any) (*http.Response, error)
	Delete(url string, headers http.Header) (*http.Response, error)
}

type client struct {
	Headers http.Header
}

func New() HttpClient {
	return &client{}
}

func (c *client) SetHeaders(headers http.Header) {
	c.Headers = headers
}

func (c *client) Get(url string, headers http.Header) (*http.Response, error) {
	return c.do(http.MethodGet, url, headers, nil)
}

func (c *client) Post(url string, headers http.Header, body any) (*http.Response, error) {
	return c.do(http.MethodPost, url, headers, body)
}

func (c *client) Put(url string, headers http.Header, body any) (*http.Response, error) {
	return c.do(http.MethodPut, url, headers, body)
}

func (c *client) Patch(url string, headers http.Header, body any) (*http.Response, error) {
	return c.do(http.MethodPatch, url, headers, body)
}

func (c *client) Delete(url string, headers http.Header) (*http.Response, error) {
	return c.do(http.MethodDelete, url, headers, nil)
}
