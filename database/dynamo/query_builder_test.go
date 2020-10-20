package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"testing"
)

func TestQueryBuilder_Select(t *testing.T) {
	test := CreateNewQupDynamo(TestConnection{})
	builder, err := test.Query("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	newBuilder := builder.Select([]string{"id", "external_id"})

	if *newBuilder.Query.ProjectionExpression != "id,external_id" {
		t.Fail()
	}
}

func TestQueryBuilder_Equals(t *testing.T) {
	test := CreateNewQupDynamo(TestConnection{})
	builder, err := test.Query("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	newBuilder := builder.Equals("external_id", 123)
	checkVal := newBuilder.Query.KeyConditions["external_id"]

	if *checkVal.AttributeValueList[0].N != "123" {
		t.Fail()
	}

	if *checkVal.ComparisonOperator != "EQ" {
		t.Fail()
	}

	if *newBuilder.Query.IndexName != "external_id-group_key" {
		t.Fail()
	}

	newBuilder = builder.Equals("non_existing_field", "Hello World")

	if len(newBuilder.Errors) == 0 {
		t.Fail()
	}

	if newBuilder.Errors[0].Error() != "unable to find index" {
		t.Fail()
	}
}

func TestQueryBuilder_Limit(t *testing.T) {
	test := CreateNewQupDynamo(TestConnection{})
	builder, err := test.Query("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	newBuilder := builder.Limit(12)

	if *newBuilder.Query.Limit != 12 {
		t.Fail()
	}
}

func TestQueryBuilder_Count(t *testing.T) {
	test := CreateNewQupDynamo(TestConnection{
		MockQueryItemCount: aws.Int64(20),
	})
	builder, err := test.Query("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	result, errs := builder.Count()

	if errs != nil {
		t.Fail()
	}

	if result != 20 {
		t.Fail()
	}
}
