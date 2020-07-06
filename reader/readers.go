package reader

import (
	"fmt"
	types "github.com/queueup-dev/qup-types"
	"io"
	"strings"
)

func NewJsonReader(stream io.Reader) *jsonReader {
	return &jsonReader{input: stream}
}

func NewXmlReader(stream io.Reader) *xmlReader {
	return &xmlReader{input: stream}
}

func NewRawReader(stream io.Reader) *rawReader {
	return &rawReader{input: stream}
}

func NewProtoReader(stream io.Reader) *protoReader {
	return &protoReader{input: stream}
}

func NewReader(contentType string, stream io.Reader) types.PayloadReader {
	contentType = strings.ToLower(contentType)

	switch {
	case strings.Contains(contentType, "text/xml"), strings.Contains(contentType, "application/xml"):
		return NewXmlReader(stream)
	case strings.Contains(contentType, "application/json"):
		return NewJsonReader(stream)
	case strings.Contains(contentType, "application/x-protobuf"), strings.Contains(contentType, "application/protobuf"):
		return NewProtoReader(stream)
	}

	return NewRawReader(stream)
}
