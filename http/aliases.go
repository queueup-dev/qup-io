package http

import "io"

func Delete(
	client Client,
	url string,
	body io.Reader,
	headers *map[string]string,
) (Result, error) {
	return Request(client, "DELETE", url, headers, body)
}

func Get(
	client Client,
	url string,
	headers *map[string]string,
) (Result, error) {
	return Request(client, "GET", url, headers, nil)
}

func Options(
	client Client,
	url string,
) (Result, error) {
	return Request(client, "OPTIONS", url, nil, nil)
}

func Patch(
	client Client,
	url string,
	body io.Reader,
	headers *map[string]string,
) (Result, error) {
	return Request(client, "PATCH", url, headers, body)
}

func Post(
	client Client,
	url string,
	body io.Reader,
	headers *map[string]string,
) (Result, error) {
	return Request(client, "POST", url, headers, body)
}
