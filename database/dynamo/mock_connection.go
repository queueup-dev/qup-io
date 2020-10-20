package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type TestConnection struct {
	MockQueryResult []map[string]*dynamodb.AttributeValue
	MockQueryError  error

	MockGetItemResult  map[string]*dynamodb.AttributeValue
	MockGetItemError   error
	MockQueryItemCount *int64

	DeleteItemError error

	MockScanResult []map[string]*dynamodb.AttributeValue
	MockScanError  error

	MockTransactError error
}

func (t TestConnection) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return &dynamodb.QueryOutput{
		Count: t.MockQueryItemCount,
		Items: t.MockQueryResult,
	}, t.MockQueryError
}

func (t TestConnection) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &dynamodb.GetItemOutput{
		Item: t.MockGetItemResult,
	}, t.MockGetItemError
}

func (t TestConnection) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return &dynamodb.DeleteItemOutput{}, t.DeleteItemError
}

func (t TestConnection) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return &dynamodb.ScanOutput{
		Count: aws.Int64(int64(len(t.MockScanResult))),
		Items: t.MockScanResult,
	}, t.MockScanError
}

func (t TestConnection) TransactWriteItems(input *dynamodb.TransactWriteItemsInput) (*dynamodb.TransactWriteItemsOutput, error) {
	return &dynamodb.TransactWriteItemsOutput{}, t.MockTransactError
}
