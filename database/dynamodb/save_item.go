package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	types "github.com/queueup-dev/qup-types"
)

func SaveItem(connection *dynamodb.DynamoDB, table string, writer types.PayloadWriter) (types.PayloadReader, error) {

	content, err := writer.Marshal()
	record := content.(map[string]*dynamodb.AttributeValue)

	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:         record,
		TableName:    aws.String(table),
		ReturnValues: aws.String("ALL_NEW"),
	}

	result, err := connection.PutItem(input)

	if err != nil {
		return nil, err
	}

	return NewDynamoReader(result.Attributes), err
}
