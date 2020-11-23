package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"testing"
)

type ExampleRecord struct {
	Id         string `dynamo:"id,key"`
	ExternalId string `dynamo:"external_id,idx|external_id-group_key"`
	GroupKey   string `dynamo:"group_key,idx|external_id-group_key"`
	OtherField string `dynamo:"other_field,idx|other_field-index"`
}

type DynamoRecordDefinition []string

func TestQupDynamo_QueryDynamoError(t *testing.T) {
	connection := TestConnection{
		MockQueryResult: []map[string]*dynamodb.AttributeValue{
			{
				"id": {
					S: aws.String("test123"),
				},
				"external_id": {
					S: aws.String("external_id123"),
				},
				"group_key": {
					S: aws.String("group_key123"),
				},
			},
		},
		MockQueryError: fmt.Errorf("unexpected Dynamo error"),
	}

	test := CreateNewQupDynamo(connection)
	builder, err := test.Query("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	_, errs := builder.Execute()

	if errs == nil || len(*errs) > 1 {
		t.Fail()
	}
}

func TestQupDynamo_QueryHappyFlow(t *testing.T) {
	connection := TestConnection{
		MockQueryResult: []map[string]*dynamodb.AttributeValue{
			{
				"id": {
					S: aws.String("test123"),
				},
				"external_id": {
					S: aws.String("external_id123"),
				},
				"group_key": {
					S: aws.String("group_key123"),
				},
			},
		},
		MockQueryError: nil,
	}

	test := CreateNewQupDynamo(connection)
	builder, err := test.Query("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	result, errs := builder.Execute()

	if errs != nil {
		t.Fail()
	}

	firstRecord := ExampleRecord{}
	err = result.First(&firstRecord)

	if err != nil {
		t.Fail()
	}

	if firstRecord.ExternalId != "external_id123" {
		t.Fail()
	}

	if firstRecord.GroupKey != "group_key123" {
		t.Fail()
	}

	if firstRecord.Id != "test123" {
		t.Fail()
	}
}

func TestQupDynamo_Retrieve(t *testing.T) {
	connection := TestConnection{
		MockGetItemResult: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("test123"),
			},
			"external_id": {
				S: aws.String("external_id123"),
			},
			"group_key": {
				S: aws.String("group_key123"),
			},
		},
		MockGetItemError: nil,
	}

	test := CreateNewQupDynamo(connection)
	record := ExampleRecord{}
	err := test.Retrieve("mockTable", "test123", &record)

	if err != nil {
		fmt.Print(err)
		t.Fail()
	}

	if record.Id != "test123" {
		t.Fail()
	}

	if record.ExternalId != "external_id123" {
		t.Fail()
	}

	if record.GroupKey != "group_key123" {
		t.Fail()
	}
}

func TestQupDynamo_RetrieveFail(t *testing.T) {
	connection := TestConnection{
		MockGetItemResult: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("test123"),
			},
			"external_id": {
				S: aws.String("external_id123"),
			},
			"group_key": {
				S: aws.String("group_key123"),
			},
		},
		MockGetItemError: fmt.Errorf("test error"),
	}

	test := CreateNewQupDynamo(connection)
	record := ExampleRecord{}
	err := test.Retrieve("mockTable", "test123", &record)

	if err == nil {
		t.Fail()
	}

	if err.Error() != "test error" {
		t.Fail()
	}
}

func TestQupDynamo_Delete(t *testing.T) {
	connection := TestConnection{
		DeleteItemError: nil,
	}

	test := CreateNewQupDynamo(connection)
	err := test.Delete("mockTable", "test123", &ExampleRecord{})

	if err != nil {
		t.Fail()
	}
}

func TestQupDynamo_DeleteFails(t *testing.T) {
	connection := TestConnection{
		DeleteItemError: fmt.Errorf("test delete error"),
	}

	test := CreateNewQupDynamo(connection)
	err := test.Delete("mockTable", "test123", &ExampleRecord{})

	if err == nil {
		t.Fail()
	}

	if err.Error() != "test delete error" {
		t.Fail()
	}
}

func TestQupDynamo_ScanFails(t *testing.T) {
	connection := TestConnection{
		MockScanResult: []map[string]*dynamodb.AttributeValue{
			{
				"id": {
					S: aws.String("test123"),
				},
				"external_id": {
					S: aws.String("external_id123"),
				},
				"group_key": {
					S: aws.String("group_key123"),
				},
			},
			{
				"id": {
					S: aws.String("test456"),
				},
				"external_id": {
					S: aws.String("external_id456"),
				},
				"group_key": {
					S: aws.String("group_key456"),
				},
			},
		},
		MockScanError: fmt.Errorf("scan fails"),
	}

	test := CreateNewQupDynamo(connection)

	target := make([]ExampleRecord, 2)
	err := test.Scan("mockTable", &target, 2)

	if err == nil {
		t.Fail()
	}

	if err.Error() != "scan fails" {
		t.Fail()
	}
}

func TestQupDynamo_Scan(t *testing.T) {
	connection := TestConnection{
		MockScanResult: []map[string]*dynamodb.AttributeValue{
			{
				"id": {
					S: aws.String("test123"),
				},
				"external_id": {
					S: aws.String("external_id123"),
				},
				"group_key": {
					S: aws.String("group_key123"),
				},
			},
			{
				"id": {
					S: aws.String("test456"),
				},
				"external_id": {
					S: aws.String("external_id456"),
				},
				"group_key": {
					S: aws.String("group_key456"),
				},
			},
		},
	}

	test := CreateNewQupDynamo(connection)

	target := make([]ExampleRecord, 2)
	err := test.Scan("mockTable", &target, 2)

	if err != nil {
		t.Fail()
	}

	if len(target) != 2 {
		t.Fail()
	}

	if target[0].GroupKey != "group_key123" {
		t.Fail()
	}

	if target[0].ExternalId != "external_id123" {
		t.Fail()
	}

	if target[0].Id != "test123" {
		t.Fail()
	}

	if target[1].GroupKey != "group_key456" {
		t.Fail()
	}

	if target[1].ExternalId != "external_id456" {
		t.Fail()
	}

	if target[1].Id != "test456" {
		t.Fail()
	}
}

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

	MockSaveError error
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
