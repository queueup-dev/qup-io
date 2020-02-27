package reader

import (
	"encoding/xml"
	"io"
	"io/ioutil"
)

type XmlReader struct {
	output io.Reader
}

func (x XmlReader) Unmarshal(object interface{}) error {
	return xml.NewDecoder(x.output).Decode(object)
}

func (x XmlReader) Valid() bool {
	var test string
	return x.Unmarshal(&test) == nil
}

func (x XmlReader) Reader() (io.Reader, error) {
	return x.output, nil
}

func (x XmlReader) Bytes() ([]byte, error) {
	reader, _ := x.Reader()

	return ioutil.ReadAll(reader)
}

func (x XmlReader) ToString() (*string, error) {
	result, err := x.Bytes()

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}
