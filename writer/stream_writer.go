package writer

import (
	"fmt"
	"io"
	"io/ioutil"
)

type StreamWriter struct {
	input io.Reader
}

func (s StreamWriter) Bytes() ([]byte, error) {
	return ioutil.ReadAll(s.input)
}

func (s StreamWriter) Marshal() (interface{}, error) {
	return fmt.Sprint(s.input), nil
}

func (s StreamWriter) Valid() bool {
	return s.input != nil
}

func (s StreamWriter) Reader() (io.Reader, error) {
	return s.input, nil
}

func (s StreamWriter) ToString() (*string, error) {
	content, err := s.Marshal()
	bytes := content.([]byte)

	if err != nil {
		return nil, err
	}

	output := string(bytes)

	return &output, nil
}

func (s StreamWriter) ContentType() string {
	return "application/octet-stream"
}
