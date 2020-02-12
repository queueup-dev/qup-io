package Dynamodb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateConnection(initializedSession *session.Session) *dynamodb.DynamoDB {
	return dynamodb.New(initializedSession)
}
