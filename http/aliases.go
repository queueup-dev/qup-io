package http

import (
	types "github.com/queueup-dev/qup-types"
)

func Delete(
	client Client,
	body types.PayloadWriter,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "DELETE", headers, body)
}

func Get(
	client Client,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "GET", headers, nil)
}

func Head(
	client Client,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "HEAD", headers, nil)
}

func Options(
	client Client,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "OPTIONS", nil, nil)
}

func Patch(
	client Client,
	body types.PayloadWriter,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "PATCH", headers, body)
}

func Post(
	client Client,
	body types.PayloadWriter,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "POST", headers, body)
}

func Put(
	client Client,
	body types.PayloadWriter,
	headers *map[string]string,
) (types.PayloadReader, HttpError, error) {
	return Request(client, "PUT", headers, body)
}
