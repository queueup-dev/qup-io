package reader

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type jsonReader struct {
	input io.Reader
}

func (j jsonReader) Unmarshal(object interface{}) error {
	return json.NewDecoder(j.input).Decode(object)
}

func (j jsonReader) Valid() bool {
	var test string
	return j.Unmarshal(&test) == nil
}

func (j jsonReader) Reader() (io.Reader, error) {
	return j.input, nil
}

func (j jsonReader) Bytes() ([]byte, error) {
	reader, _ := j.Reader()

	return ioutil.ReadAll(reader)
}

func (j jsonReader) ToString() (*string, error) {
	result, err := j.Bytes()

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}

func (j jsonReader) ContentType() string {
	return "application/json"
}
