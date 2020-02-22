package http

import (
	"encoding/xml"
	"io"
)

type XmlResult struct {
	output io.Reader
}

func (x XmlResult) Unmarshal(object interface{}) error {
	return xml.NewDecoder(x.output).Decode(object)
}
