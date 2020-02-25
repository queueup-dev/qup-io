package http

import (
	"github.com/queueup-dev/qup-io/reader"
	"github.com/queueup-dev/qup-types"
	"net/http"
)

func Request(
	client Client,
	method string,
	headers *map[string]string,
	body types.PayloadWriter,
) (types.PayloadReader, HttpError, error) {
	bodyReader, err := body.Reader()
	url := body.GetTarget()

	if err != nil {
		return nil, nil, err
	}

	request, err := http.NewRequest(method, url, bodyReader)

	if err != nil {
		return nil, nil, err
	}

	for key, value := range *headers {
		request.Header.Add(key, value)
	}

	response, err := client.Do(request)

	if err != nil {
		return nil, nil, err
	}

	result, err := createResponseObject(response, url)

	if err != nil {
		return nil, nil, err
	}

	if !isSuccessful(response) {
		return nil, NewHttpError(response.StatusCode, result), nil
	}

	return result, nil, nil
}

func createResponseObject(response *http.Response, url string) (types.PayloadReader, error) {

	contentType := response.Header.Get("Content-Type")

	switch contentType {
	case "application/xml", "text/xml":
		return reader.NewXmlReader(response.Body, url), nil
	case "application/json", "application/json+error":
		return reader.NewJsonReader(response.Body, url), nil
	}

	return reader.NewRawReader(response.Body, url), nil
}

func isSuccessful(response *http.Response) bool {
	if response.StatusCode >= 400 {
		return false
	}
	return true
}
