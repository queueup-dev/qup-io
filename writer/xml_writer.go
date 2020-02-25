package writer

import (
	"encoding/xml"
	"io"
	"strings"
)

type XmlWriter struct {
	input interface{}
}

func (x XmlWriter) Marshal() ([]byte, error) {
	return xml.Marshal(x.input)
}

func (x XmlWriter) Valid() bool {
	_, err := x.Marshal()

	if err != nil {
		return false
	}

	return true
}

func (x XmlWriter) Reader() (io.Reader, error) {
	content, err := x.Marshal()

	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(content)), nil
}

func (x XmlWriter) ToString() (*string, error) {
	content, err := x.Marshal()

	if err != nil {
		return nil, err
	}

	output := string(content)
	return &output, nil
}
