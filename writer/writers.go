package writer

import (
	"github.com/golang/protobuf/proto"
	"io"
)

func NewJsonWriter(input interface{}) *JsonWriter {
	return &JsonWriter{input: input}
}

func NewXmlWriter(input interface{}) *XmlWriter {
	return &XmlWriter{input: input}
}

func NewRawWriter(input interface{}) *RawWriter {
	return &RawWriter{input: input}
}

func NewProtoWriter(input proto.Message) *ProtoWriter {
	return &ProtoWriter{input: input}
}

func NewFormEncodeWriter(input interface{}) *FormEncodeWriter {
	return &FormEncodeWriter{input: input}
}

func NewStreamWriter(input io.Reader) *StreamWriter {
	return &StreamWriter{input: input}
}
