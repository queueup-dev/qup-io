package reader

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type JsonReader struct {
	output io.Reader
}

func (j JsonReader) Unmarshal(object interface{}) error {
	return json.NewDecoder(j.output).Decode(object)
}

func (j JsonReader) Valid() bool {
	var test string
	return j.Unmarshal(&test) == nil
}

func (j JsonReader) Reader() (io.Reader, error) {
	return j.output, nil
}

func (j JsonReader) Bytes() ([]byte, error) {
	reader, _ := j.Reader()

	return ioutil.ReadAll(reader)
}

func (j JsonReader) ToString() (*string, error) {
	result, err := j.Bytes()

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}