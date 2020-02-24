package http

import (
	"net/http"
	"time"
)

type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

func NewDefaultClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}
