package dynamodb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"io"
)

type DynamoWriter struct {
	input interface{}
}

func (d DynamoWriter) Bytes() ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(d.input)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (d DynamoWriter) Marshal() (interface{}, error) {
	return dynamodbattribute.MarshalMap(d.input)
}

func (d DynamoWriter) Valid() bool {
	_, err := d.Marshal()

	if err != nil {
		return false
	}

	return true
}

func (d DynamoWriter) Reader() (io.Reader, error) {
	byteOutput, err := d.Bytes()

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(byteOutput), nil
}

func (d DynamoWriter) ToString() (*string, error) {
	content, err := d.Marshal()

	if err != nil {
		return nil, err
	}

	output := fmt.Sprint(content)
	return &output, nil
}

func NewDynamoWriter(input interface{}) *DynamoWriter {
	return &DynamoWriter{input: input}
}
