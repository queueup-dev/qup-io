package http

import (
	"io"
	"net/http"
)

func Request(
	client Client,
	method string,
	url string,
	headers *map[string]string,
	body io.Reader,
) (Result, HttpError, error) {
	request, err := http.NewRequest(method, url, body)

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

	result, err := createResponseObject(response)

	if err != nil {
		return nil, nil, err
	}

	if !isSuccessful(response) {
		return nil, NewHttpError(response.StatusCode, result), nil
	}

	return result, nil, nil
}

func createResponseObject(response *http.Response) (Result, error) {

	contentType := response.Header.Get("Content-Type")

	switch contentType {
	case "application/xml", "text/xml":
		return &XmlResult{
			output: response.Body,
		}, nil
	case "application/json", "application/json+error":
		return &JsonResult{
			output: response.Body,
		}, nil
	}

	return &RawResult{
		output: response.Body,
	}, nil
}

func isSuccessful(response *http.Response) bool {
	if response.StatusCode >= 400 {
		return false
	}
	return true
}
