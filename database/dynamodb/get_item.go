package dynamodb

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GetItem(connection *dynamodb.DynamoDB, table string, primaryKey string, value string) (*DynamoReader, error) {
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

	if result.Item == nil {
		return nil, errors.New("item not found")
	}

	return &DynamoReader{output: result.Item}, nil
}
