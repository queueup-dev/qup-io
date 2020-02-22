package http

import (
	"io"
)

func Post(
	client Client,
	url string,
	headers map[string]string,
	body io.Reader,
) (Result, error) {
	return Request(client, "POST", url, &headers, body)
}
