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
) (Result, error) {
	request, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	for key, value := range *headers {
		request.Header.Add(key, value)
	}

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	_, httpError := isSuccessful(response)

	if httpError != nil {
		return nil, *httpError
	}

	return createResponseObject(response)
}

func createResponseObject(response *http.Response) (Result, error) {

	contentType := response.Header.Get("Content-Type")

	switch contentType {
	case "application/xml", "text/xml":
		return &XmlResult{
			output: response.Body,
		}, nil
	}

	return &JsonResult{
		output: response.Body,
	}, nil
}

func isSuccessful(response *http.Response) (bool, *HttpError) {

	if response.StatusCode >= 400 {
		return false, &HttpError{
			statusCode: response.StatusCode,
			message:    "Invalid response status",
		}
	}

	return true, nil
}
