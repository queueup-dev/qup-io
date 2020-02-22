package Dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func DeleteItem(connection *dynamodb.DynamoDB, table string, primaryKey string, value string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			primaryKey: {
				S: aws.String(value),
			},
		},
		TableName: aws.String(table),
	}

	_, err := connection.DeleteItem(input)

	return err
}
