package http

import (
	types "github.com/queueup-dev/qup-types"
	"net/http"
	"strconv"
)

type HttpError interface {
	StatusCode() int
	Error() string
	HttpResponse() types.PayloadReader
}

type Error struct {
	statusCode int
	message    string
	response   types.PayloadReader
}

func (e Error) HttpResponse() types.PayloadReader {
	return e.response
}

func (e Error) StatusCode() int {
	return e.statusCode
}

func (e Error) Error() string {
	return e.message
}

func NewHttpError(statusCode int, response types.PayloadReader) HttpError {
	message := "[HTTP:" + strconv.Itoa(statusCode) + "] " + http.StatusText(statusCode)

	return &Error{
		statusCode: statusCode,
		message:    message,
		response:   response,
	}
}
