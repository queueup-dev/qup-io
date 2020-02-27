package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	types "github.com/queueup-dev/qup-types"
)

func SaveItem(connection dynamodb.DynamoDB, table string, writer types.PayloadWriter) error {

	content, err := writer.Marshal()
	record := content.(map[string]*dynamodb.AttributeValue)

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
