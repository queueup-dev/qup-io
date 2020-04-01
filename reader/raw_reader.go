package reader

import (
	"errors"
	"io"
	"io/ioutil"
)

type rawReader struct {
	input io.Reader
}

func (r rawReader) Unmarshal(object interface{}) error {
	return errors.New("unable to unmarshal plain text")
}

func (r rawReader) Valid() bool {
	return true
}

func (r rawReader) Reader() (io.Reader, error) {
	return r.input, nil
}

func (r rawReader) Bytes() ([]byte, error) {
	reader, _ := r.Reader()

	return ioutil.ReadAll(reader)
}

func (r rawReader) ToString() (*string, error) {
	result, err := ioutil.ReadAll(r.input)

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}

func (r rawReader) ContentType() string {
	return "text/plain"
}
