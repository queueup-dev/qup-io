package Dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ResultSet struct {
	output *dynamodb.GetItemOutput
}

func (r ResultSet) Unmarshal(object interface{}) (interface{}, error) {
	err := dynamodbattribute.UnmarshalMap(r.output.Item, object)

	if err != nil {
		return nil, err
	}

	return object, nil
}
