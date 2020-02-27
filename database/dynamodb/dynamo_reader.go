package dynamodb

import (
	"bytes"
	"encoding/gob"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"io"
)

type DynamoReader struct {
	output *dynamodb.GetItemOutput
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
	return dynamodbattribute.UnmarshalMap(d.output.Item, object)
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
