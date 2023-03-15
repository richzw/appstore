package appstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"time"
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type DoFunc func(*http.Request) (*http.Response, error)

func (d DoFunc) Do(req *http.Request) (*http.Response, error) {
	return d(req)
}

func SetHeader(c HTTPClient, key string, value ...string) DoFunc {
	key = textproto.CanonicalMIMEHeaderKey(key)
	return func(req *http.Request) (*http.Response, error) {
		if req.Header == nil {
			req.Header = make(http.Header)
		}
		req.Header[key] = value
		return c.Do(req)
	}
}

func AddHeader(client HTTPClient, key string, value ...string) DoFunc {
	key = textproto.CanonicalMIMEHeaderKey(key)
	return func(req *http.Request) (*http.Response, error) {
		if req.Header == nil {
			req.Header = make(http.Header)
		}
		req.Header[key] = append(req.Header[key], value...)
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
			return resp, fmt.Errorf("received invalid status code: %d", resp.StatusCode)
		}
		return resp, nil
	}
}

type Initializer func(HTTPClient) (DoFunc, error)

func SetInitializer(c HTTPClient, init Initializer) DoFunc {
	var f DoFunc
	return func(req *http.Request) (*http.Response, error) {
		var err error
		f, err = init(c)
		if err != nil {
			return nil, err
		}

		return f.Do(req)
	}
}

func SetRequest(ctx context.Context, c HTTPClient, method string, url string) DoFunc {
	return func(_ *http.Request) (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, err
		}
		return c.Do(req)
	}
}

type Marshaller func(v any) ([]byte, error)
type Unmarshaller func(b []byte, v any) error

func SetRequestBody(c HTTPClient, m Marshaller, v any) DoFunc {
	return func(req *http.Request) (*http.Response, error) {
		if m == nil {
			switch t := v.(type) {
			case []byte:
				req.Body = io.NopCloser(bytes.NewReader(t))
			case io.ReadCloser:
				req.Body = t
			case io.Reader:
				req.Body = io.NopCloser(t)
			default:
				return nil, fmt.Errorf("failed to marshal body type: %v", v)
			}
		} else {
			b, err := m(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			req.Body = io.NopCloser(bytes.NewReader(b))
		}
		return c.Do(req)
	}
}

func SetRequestBodyJSON(c HTTPClient, v any) DoFunc {
	c = SetHeader(c, "Content-Type", "application/json")
	return SetRequestBody(c, json.Marshal, v)
}

func SetResponseBodyHandler(c HTTPClient, u Unmarshaller, ptr any) DoFunc {
	c = RequireResponseBody(c)
	return func(req *http.Request) (*http.Response, error) {
		resp, err := c.Do(req)
		if err != nil {
			return resp, err
		}
		b, err := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()
		if err != nil {
			return resp, err
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(b))
		if err = u(b, ptr); err != nil {
			return resp, err
		}
		if closeErr != nil {
			return resp, closeErr
		}
		return resp, nil
	}
}

func RateLimit(c HTTPClient, reqPerMin int) DoFunc {
	ticker := time.NewTicker(time.Second * time.Duration(60))
	ch := make(chan struct{}, reqPerMin)
	go func() {
		for range ticker.C {
			for i := 0; i < reqPerMin; i++ {
				select {
				case ch <- struct{}{}:
				default:
					break
				}
			}
		}
	}()
	return func(req *http.Request) (*http.Response, error) {
		select {
		case <-ch:
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
		return c.Do(req)
	}
}

func ShouldRetryDefault(status int, err error) bool {
	if 500 <= status && status <= 599 {
		return true
	}
	if status == http.StatusTooManyRequests {
		return true
	}
	if err == io.ErrUnexpectedEOF {
		return true
	}

	// If error unwrapping is available, use this to examine wrapped errors.
	if err, ok := err.(interface{ Unwrap() error }); ok {
		return ShouldRetryDefault(status, err.Unwrap())
	}
	return false
}

func SetRetry(c HTTPClient, bo Backoff, shouldRetry func(int, error) bool) DoFunc {
	return func(req *http.Request) (*http.Response, error) {
		var resp *http.Response
		var err error
		var pause time.Duration

		for {
			select {
			case <-req.Context().Done():
				if err == nil {
					err = req.Context().Err()
				}
				return resp, err
			case <-time.After(pause):
			}

			resp, err = c.Do(req)

			var status int
			if resp != nil {
				status = resp.StatusCode
			}

			if req.GetBody == nil || !shouldRetry(status, err) {
				break
			}
			var errBody error
			req.Body, errBody = req.GetBody()
			if errBody != nil {
				break
			}

			pause = bo.Pause()
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}

			if pause < 0 {
				return nil, fmt.Errorf("sendAndRetry timeout for url %s", req.URL)
			}
		}
		return resp, err
	}
}
