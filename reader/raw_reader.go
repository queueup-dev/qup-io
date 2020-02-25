package reader

import (
	"errors"
	"io"
	"io/ioutil"
)

type RawReader struct {
	output io.Reader
	source string
}

func (r RawReader) GetSource() string {
	return r.source
}

func (r RawReader) Unmarshal(object interface{}) error {
	return errors.New("unable to unmarshal plain text")
}

func (r RawReader) Valid() bool {
	return true
}

func (r RawReader) Reader() io.Reader {
	return r.output
}

func (r RawReader) Bytes() ([]byte, error) {
	return ioutil.ReadAll(r.Reader())
}

func (r RawReader) ToString() (*string, error) {
	result, err := ioutil.ReadAll(r.output)

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}
