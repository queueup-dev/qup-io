package http

import (
	"io"
)

func Patch(
	client Client,
	url string,
	body io.Reader,
	headers *map[string]string,
) (Result, error) {
	return Request(client, "PATCH", url, headers, body)
}
