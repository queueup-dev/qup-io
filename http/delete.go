package http

import (
	"io"
)

func Delete(
	client Client,
	url string,
	body io.Reader,
	headers *map[string]string,
) (Result, error) {
	return Request(client, "DELETE", url, headers, body)
}
