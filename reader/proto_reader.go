package reader

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
)

type protoReader struct {
	input io.Reader
}

func (p protoReader) Unmarshal(object interface{}) error {

	protoMessage, ok := object.(proto.Message)

	if !ok {
		return errors.New("only proto-messages are supported")
	}

	result, err := p.Bytes()

	if err != nil {
		return err
	}

	return proto.Unmarshal(result, protoMessage)
}

func (p protoReader) Valid() bool {
	return true
}

func (p protoReader) Reader() (io.Reader, error) {
	return p.input, nil
}

func (p protoReader) Bytes() ([]byte, error) {
	reader, _ := p.Reader()

	return ioutil.ReadAll(reader)
}

func (p protoReader) ToString() (*string, error) {
	result, err := p.Bytes()

	if err != nil {
		return nil, err
	}

	stringResult := string(result)
	return &stringResult, nil
}

func (p protoReader) ContentType() string {
	return "application/proto"
}
