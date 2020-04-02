package dynamodb

import (
	"bytes"
	"encoding/gob"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"io"
)

type DynamoReader struct {
	output map[string]*dynamodb.AttributeValue
}

func (d DynamoReader) Bytes() ([]byte, error) {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(d.output)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (d DynamoReader) Unmarshal(object interface{}) error {
	return dynamodbattribute.UnmarshalMap(d.output, object)
}

func (d DynamoReader) Valid() bool {
	var test string
	return d.Unmarshal(&test) == nil
}

func (d DynamoReader) Reader() (io.Reader, error) {
	outputBytes, err := d.Bytes()

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(outputBytes), nil
}

func (d DynamoReader) ContentType() string {
	return "application/dynamodb"
}

func (d DynamoReader) ToString() (*string, error) {
	return nil, nil
}

func NewDynamoReader(input map[string]*dynamodb.AttributeValue) *DynamoReader {
	return &DynamoReader{output: input}
}
