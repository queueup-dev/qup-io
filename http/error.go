package http

type HttpError struct {
	statusCode int
	message    string
}

func (e HttpError) StatusCode() int {
	return e.statusCode
}

func (e HttpError) Error() string {
	return e.message
}
