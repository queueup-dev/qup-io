package http

import (
	"encoding/xml"
	"io"
	"io/ioutil"
)

type XmlResult struct {
	output io.Reader
}

func (x XmlResult) Unmarshal(object interface{}) error {
	return xml.NewDecoder(x.output).Decode(object)
}

func (x XmlResult) Valid() bool {
	var test string
	return x.Unmarshal(&test) == nil
}

func (x XmlResult) ToString() (*string, error) {
	result, err := ioutil.ReadAll(x.output)

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}
