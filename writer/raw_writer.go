package writer

import (
	"fmt"
	"io"
	"strings"
)

type RawWriter struct {
	input interface{}
}

func (r RawWriter) Bytes() ([]byte, error) {
	marshalOutput, err := r.Marshal()

	if err != nil {
		return nil, err
	}

	return marshalOutput.([]byte), nil
}

func (r RawWriter) Marshal() (interface{}, error) {
	return []byte(fmt.Sprint(r.input)), nil
}

func (r RawWriter) Valid() bool {
	_, err := r.Marshal()

	if err != nil {
		return false
	}

	return true
}

func (r RawWriter) Reader() (io.Reader, error) {
	content, err := r.Marshal()
	bytes := content.([]byte)

	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(bytes)), nil
}

func (r RawWriter) ToString() (*string, error) {
	content, err := r.Marshal()
	bytes := content.([]byte)

	if err != nil {
		return nil, err
	}

	output := string(bytes)
	return &output, nil
}
