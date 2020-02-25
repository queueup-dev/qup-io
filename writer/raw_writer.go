package writer

import (
	"fmt"
	"io"
	"strings"
)

type RawWriter struct {
	input  interface{}
	target string
}

func (r RawWriter) GetTarget() string {
	return r.target
}

func (r RawWriter) Marshal() ([]byte, error) {
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

	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(content)), nil
}

func (r RawWriter) ToString() (*string, error) {
	content, err := r.Marshal()

	if err != nil {
		return nil, err
	}

	output := string(content)
	return &output, nil
}
