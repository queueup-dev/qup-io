package reader

import (
	"encoding/xml"
	"io"
	"io/ioutil"
)

type xmlReader struct {
	input  io.Reader
	output interface{}
}

func (x xmlReader) Unmarshal(object interface{}) error {
	return xml.NewDecoder(x.input).Decode(object)
}

func (x xmlReader) Valid() bool {
	var test string
	return x.Unmarshal(&test) == nil
}

func (x xmlReader) Reader() (io.Reader, error) {
	return x.input, nil
}

func (x xmlReader) Bytes() ([]byte, error) {
	reader, _ := x.Reader()

	return ioutil.ReadAll(reader)
}

func (x xmlReader) ToString() (*string, error) {
	result, err := x.Bytes()

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}

func (x xmlReader) ContentType() string {
	return "text/xml"
}
