package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type JsonResult struct {
	output io.Reader
}

func (j JsonResult) Unmarshal(object interface{}) error {
	return json.NewDecoder(j.output).Decode(object)
}

func (j JsonResult) Valid() bool {
	var test string
	return j.Unmarshal(&test) == nil
}

func (j JsonResult) ToString() (*string, error) {
	result, err := ioutil.ReadAll(j.output)

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}
