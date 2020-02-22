package http

import (
	"io"
)

func Post(
	client Client,
	url string,
	body io.Reader,
	headers *map[string]string,
) (Result, error) {
	return Request(client, "POST", url, headers, body)
}
