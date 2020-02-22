package http

import (
	"io"
)

func Get(
	client Client,
	url string,
	headers map[string]string,
	body io.Reader,
) (Result, error) {
	return Request(client, "GET", url, &headers, body)
}
