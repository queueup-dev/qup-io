package http

import (
	"io"
)

func Delete(
	client Client,
	url string,
	headers map[string]string,
	body io.Reader,
) (Result, error) {
	return Request(client, "DELETE", url, &headers, body)
}
