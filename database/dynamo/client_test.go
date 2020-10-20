package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"testing"
)

type ExampleRecord struct {
	Id         string `dynamo:"id,pkey"`
	ExternalId string `dynamo:"external_id,gsi|external_id-group_key"`
	GroupKey   string `dynamo:"group_key,gsi|external_id-group_key"`
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
