package http

import (
	"net/http"
	"strconv"
)

type HttpError interface {
	StatusCode() int
	Error() string
	HttpResponse() Result
}

type Error struct {
	statusCode int
	message    string
	response   Result
}

func (e Error) HttpResponse() Result {
	return e.response
}

func (e Error) StatusCode() int {
	return e.statusCode
}

func (e Error) Error() string {
	return e.message
}

func NewHttpError(statusCode int, response Result) HttpError {
	message := "[HTTP:" + strconv.Itoa(statusCode) + "] " + http.StatusText(statusCode)

	return &Error{
		statusCode: statusCode,
		message:    message,
		response:   response,
	}
}
