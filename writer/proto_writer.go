package writer

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"io"
)

type ProtoWriter struct {
	input proto.Message
}

func (p ProtoWriter) Bytes() ([]byte, error) {
	marshalledInput, err := p.Marshal()

	return marshalledInput.([]byte), err
}

func (p ProtoWriter) Marshal() (interface{}, error) {
	return proto.Marshal(p.input)
}

func (p ProtoWriter) Reader() (io.Reader, error) {
	content, err := p.Bytes()

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(content), nil
}

func (p ProtoWriter) Valid() bool {
	return true
}

func (p ProtoWriter) ToString() (*string, error) {
	content, err := p.Bytes()

	if err != nil {
		return nil, err
	}

	output := string(content)
	return &output, nil
}

func (p ProtoWriter) ContentType() string {
	return "application/proto"
}
