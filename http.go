package appstore

import (
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type DoFunc func(*http.Request) (*http.Response, error)

func (d DoFunc) Do(req *http.Request) (*http.Response, error) {
	return d(req)
}

func AddHeader(client HTTPClient, key, value string) DoFunc {
	return func(req *http.Request) (*http.Response, error) {
		if req.Header == nil {
			req.Header = make(http.Header)
		}
		req.Header.Set(key, value)
		return client.Do(req)
	}
}

func RequireResponseBody(c HTTPClient) DoFunc {
	return func(req *http.Request) (*http.Response, error) {
		resp, err := c.Do(req)
		if err != nil {
			return resp, err
		}
		if resp.Body == nil {
			return resp, fmt.Errorf("nil response body")
		}
		return resp, nil
	}
}

func RequireResponseStatus(c HTTPClient, status ...int) DoFunc {
	if len(status) == 0 {
		status = []int{http.StatusOK}
	}
	valid := make(map[int]bool, len(status))
	for _, s := range status {
		valid[s] = true
	}
	return func(req *http.Request) (*http.Response, error) {
		resp, err := c.Do(req)
		if err != nil {
			return resp, err
		}
		if !valid[resp.StatusCode] {
			return resp, fmt.Errorf("received invalid satus code: %d", resp.StatusCode)
		}
		return resp, nil
	}
}
