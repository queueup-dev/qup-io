package reader

import (
	types "github.com/queueup-dev/qup-types"
	"io"
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
	switch contentType {
	case "application/xml", "text/xml":
		return NewXmlReader(stream)
	case "application/json", "application/problem+json":
		return NewJsonReader(stream)
	case "application/protobuf", "application/x-protobuf":
		return NewProtoReader(stream)
	}

	return NewRawReader(stream)
}
