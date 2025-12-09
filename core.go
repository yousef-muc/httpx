package httpx

import (
	"errors"
	"net/http"
)

func (c *client) do(method, url string, headers http.Header, body any) (*http.Response, error) {

	_ = headers
	_ = body

	client := &http.Client{}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.New("unable to create a new request")
	}

	return client.Do(request)

}
