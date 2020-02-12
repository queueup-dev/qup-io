package Dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GetItem(connection *dynamodb.DynamoDB, table string, primaryKey string, value string) (*ResultSet, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			primaryKey: {
				S: aws.String(value),
			},
		},
	}

	result, err := connection.GetItem(input)

	if err != nil {
		return nil, err
	}

	return &ResultSet{output: result}, nil
}
