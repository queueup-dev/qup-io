package writer

import (
	"encoding/json"
	"io"
	"strings"
)

type JsonWriter struct {
	input interface{}
}

func (j JsonWriter) Marshal() ([]byte, error) {
	return json.Marshal(j.input)
}

func (j JsonWriter) Valid() bool {
	_, err := j.Marshal()

	if err != nil {
		return false
	}

	return true
}

func (j JsonWriter) Reader() (io.Reader, error) {
	content, err := j.Marshal()

	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(content)), nil
}

func (j JsonWriter) ToString() (*string, error) {
	content, err := j.Marshal()

	if err != nil {
		return nil, err
	}

	output := string(content)
	return &output, nil
}
