package http

import (
	types "github.com/queueup-dev/qup-types"
)

func Delete(
	client Client,
	url string,
	body types.PayloadWriter,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "DELETE", url, headers, body)
}

func Get(
	client Client,
	url string,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "GET", url, headers, nil)
}

func Head(
	client Client,
	url string,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "HEAD", url, headers, nil)
}

func Options(
	client Client,
	url string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "OPTIONS", url, nil, nil)
}

func Patch(
	client Client,
	url string,
	body types.PayloadWriter,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "PATCH", url, headers, body)
}

func Post(
	client Client,
	url string,
	body types.PayloadWriter,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "POST", url, headers, body)
}

func Put(
	client Client,
	url string,
	body types.PayloadWriter,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "PUT", url, headers, body)
}
