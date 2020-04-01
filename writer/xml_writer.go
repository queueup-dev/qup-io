package writer

import (
	"encoding/xml"
	"io"
	"strings"
)

type XmlWriter struct {
	input interface{}
}

func (x XmlWriter) Bytes() ([]byte, error) {
	marshalOutput, err := x.Marshal()

	if err != nil {
		return nil, err
	}

	return marshalOutput.([]byte), nil
}

func (x XmlWriter) Marshal() (interface{}, error) {
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
	bytes := content.([]byte)

	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(bytes)), nil
}

func (x XmlWriter) ToString() (*string, error) {
	content, err := x.Marshal()
	bytes := content.([]byte)

	if err != nil {
		return nil, err
	}

	output := string(bytes)
	return &output, nil
}

func (x XmlWriter) ContentType() string {
	return "text/xml"
}
