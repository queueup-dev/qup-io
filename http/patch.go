package http

import (
	"io"
)

func Patch(
	client Client,
	url string,
	headers map[string]string,
	body io.Reader,
) (Result, error) {
	return Request(client, "PATCH", url, &headers, body)
}
