package reader

import (
	"errors"
	"io"
	"io/ioutil"
)

type RawReader struct {
	output io.Reader
}

func (r RawReader) Unmarshal(object interface{}) error {
	return errors.New("unable to unmarshal plain text")
}

func (r RawReader) Valid() bool {
	return true
}

func (r RawReader) Reader() (io.Reader, error) {
	return r.output, nil
}

func (r RawReader) Bytes() ([]byte, error) {
	reader, _ := r.Reader()

	return ioutil.ReadAll(reader)
}

func (r RawReader) ToString() (*string, error) {
	result, err := ioutil.ReadAll(r.output)

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}
