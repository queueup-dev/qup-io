package writer

import (
	"encoding/json"
	"io"
	"strings"
)

type JsonWriter struct {
	input interface{}
}

func (j JsonWriter) Bytes() ([]byte, error) {
	marshalOutput, err := j.Marshal()

	if err != nil {
		return nil, err
	}

	return marshalOutput.([]byte), nil
}

func (j JsonWriter) Marshal() (interface{}, error) {
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
	bytes := content.([]byte)

	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(bytes)), nil
}

func (j JsonWriter) ToString() (*string, error) {
	content, err := j.Marshal()
	bytes := content.([]byte)

	if err != nil {
		return nil, err
	}

	output := string(bytes)
	return &output, nil
}

func (j JsonWriter) ContentType() string {
	return "application/json"
}
