package http

import (
	"errors"
	"io"
	"io/ioutil"
)

type RawResult struct {
	output io.Reader
}

func (r RawResult) Unmarshal(object interface{}) error {
	return errors.New("unable to unmarshal plain text")
}

func (r RawResult) Valid() bool {
	return true
}

func (r RawResult) ToString() (*string, error) {
	result, err := ioutil.ReadAll(r.output)

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}
