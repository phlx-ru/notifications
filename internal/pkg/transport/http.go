package transport

import (
	"net"
	"net/http"
	"time"
)

const (
	defaultTimout              = 30 * time.Second
	defaultMaxIdleConnsPerHost = 255
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultTimout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConnsPerHost:   defaultMaxIdleConnsPerHost,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}
