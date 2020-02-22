package http

import (
	"encoding/json"
	"io"
)

type JsonResult struct {
	output io.Reader
}

func (j JsonResult) Unmarshal(object interface{}) error {
	return json.NewDecoder(j.output).Decode(object)
}
