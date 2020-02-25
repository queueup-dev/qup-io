package reader

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type JsonReader struct {
	output io.Reader
	source string
}

func (j JsonReader) GetSource() string {
	return j.source
}

func (j JsonReader) Unmarshal(object interface{}) error {
	return json.NewDecoder(j.output).Decode(object)
}

func (j JsonReader) Valid() bool {
	var test string
	return j.Unmarshal(&test) == nil
}

func (j JsonReader) Reader() io.Reader {
	return j.output
}

func (j JsonReader) Bytes() ([]byte, error) {
	return ioutil.ReadAll(j.output)
}

func (j JsonReader) ToString() (*string, error) {
	result, err := j.Bytes()

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}
