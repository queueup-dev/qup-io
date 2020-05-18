package http

import (
	"github.com/queueup-dev/qup-io/reader"
	"github.com/queueup-dev/qup-types"
	"io"
	"net/http"
	"strings"
)

type Headers map[string]string

func Request(
	client Client,
	method string,
	url string,
	headers *Headers,
	body types.PayloadWriter,
) (types.PayloadReader, Headers, HttpError, error) {
	var bodyReader io.Reader
	var err error

	if body == nil {
		bodyReader = nil
	} else {
		bodyReader, err = body.Reader()
	}

	if err != nil {
		return nil, nil, nil, err
	}

	request, err := http.NewRequest(method, url, bodyReader)

	if err != nil {
		return nil, nil, nil, err
	}

	if headers != nil {
		for key, value := range *headers {
			request.Header.Add(key, value)
		}
	}

	if body != nil {
		if _, ok := request.Header["Content-Type"]; !ok {
			request.Header.Add("Content-Type", body.ContentType())
		}
	}

	response, err := client.Do(request)

	if err != nil {
		return nil, nil, nil, err
	}

	result, err := createResponseObject(response)

	if err != nil {
		return nil, nil, nil, err
	}

	if !isSuccessful(response) {
		return result, simplifyHeaders(response.Header), NewHttpError(response.StatusCode, result), nil
	}

	return result, simplifyHeaders(response.Header), nil, nil
}

func simplifyHeaders(headers http.Header) Headers {
	simplifiedHeaders := Headers{}

	for key, values := range headers {
		simplifiedHeaders[key] = strings.Join(values, " ")
	}

	return simplifiedHeaders
}

func createResponseObject(response *http.Response) (types.PayloadReader, error) {
	return reader.NewReader(response.Header.Get("Content-Type"), response.Body), nil
}

func isSuccessful(response *http.Response) bool {
	if response.StatusCode >= 400 {
		return false
	}
	return true
}
