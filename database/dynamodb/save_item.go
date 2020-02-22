package Dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func SaveItem(connection dynamodb.DynamoDB, table string, data interface{}) error {
	record, err := dynamodbattribute.MarshalMap(data)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      record,
		TableName: aws.String(table),
	}

	_, err = connection.PutItem(input)

	return err
}
