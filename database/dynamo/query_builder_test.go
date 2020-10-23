package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

func TestQueryBuilder_ChangingIndexFails(t *testing.T) {
	test := CreateNewQupDynamo(TestConnection{})
	builder, err := test.Query("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	newBuilder := builder.GreaterThan("external_id", 123)
	newBuilder = newBuilder.GreaterThan("other_field", 456)

	_, errs := newBuilder.Execute()

	if errs == nil || len(*errs) != 1 {
		t.Fail()
	}

	newBuilder = builder.GreaterThan("external_id", 123)
	newBuilder = newBuilder.GreaterThan("other_field", 456)

	_, errs = newBuilder.Count()

	if errs == nil || len(*errs) != 1 {
		t.Fail()
	}
}

func TestQueryBuilder_Execute(t *testing.T) {
	test := CreateNewQupDynamo(TestConnection{
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
	})

	builder, _ := test.Query("mockTable", ExampleRecord{})

	builder.Errors = append(builder.Errors, fmt.Errorf("test error"))

	_, errs := builder.Execute()

	if errs == nil {
		t.Fail()
	}

	builder.Errors = nil

	res, errs := builder.Execute()

	if errs != nil {
		t.Fail()
	}

	record := ExampleRecord{}
	err := res.First(&record)

	if err != nil {
		t.Fail()
	}

	if record.ExternalId != "external_id123" {
		t.Fail()
	}

	if record.GroupKey != "group_key123" {
		t.Fail()
	}

	if record.Id != "test123" {
		t.Fail()
	}

	builder.Connection = TestConnection{
		MockQueryError: fmt.Errorf("failed due to an unkown curse"),
	}

	_, errs = builder.Execute()

	if errs == nil || len(*errs) > 1 {
		t.Fail()
	}

	if (*errs)[0].Error() != "failed due to an unkown curse" {
		t.Fail()
	}
}

func TestQueryBuilder_GreaterThan(t *testing.T) {
	test := CreateNewQupDynamo(TestConnection{})
	builder, err := test.Query("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	newBuilder := builder.GreaterThan("external_id", 123)
	checkVal := newBuilder.Query.KeyConditions["external_id"]

	if *checkVal.AttributeValueList[0].N != "123" {
		t.Fail()
	}

	if *checkVal.ComparisonOperator != "GT" {
		t.Fail()
	}

	if *newBuilder.Query.IndexName != "external_id-group_key" {
		t.Fail()
	}

	newBuilder = builder.GreaterThan("non_existing_field", "Hello World")

	if len(newBuilder.Errors) == 0 {
		t.Fail()
	}

	if newBuilder.Errors[0].Error() != "unable to find index" {
		t.Fail()
	}
}

func TestQueryBuilder_EqualOrLowerThan(t *testing.T) {
	test := CreateNewQupDynamo(TestConnection{})
	builder, err := test.Query("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	newBuilder := builder.EqualOrLowerThan("external_id", 123)
	checkVal := newBuilder.Query.KeyConditions["external_id"]

	if *checkVal.AttributeValueList[0].N != "123" {
		t.Fail()
	}

	if *checkVal.ComparisonOperator != "LE" {
		t.Fail()
	}

	if *newBuilder.Query.IndexName != "external_id-group_key" {
		t.Fail()
	}

	newBuilder = builder.EqualOrLowerThan("non_existing_field", "Hello World")

	if len(newBuilder.Errors) == 0 {
		t.Fail()
	}

	if newBuilder.Errors[0].Error() != "unable to find index" {
		t.Fail()
	}
}

func TestQueryBuilder_EqualOrGreaterThan(t *testing.T) {
	test := CreateNewQupDynamo(TestConnection{})
	builder, err := test.Query("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	newBuilder := builder.EqualOrGreaterThan("external_id", 123)
	checkVal := newBuilder.Query.KeyConditions["external_id"]

	if *checkVal.AttributeValueList[0].N != "123" {
		t.Fail()
	}

	if *checkVal.ComparisonOperator != "GE" {
		t.Fail()
	}

	if *newBuilder.Query.IndexName != "external_id-group_key" {
		t.Fail()
	}

	newBuilder = builder.EqualOrGreaterThan("non_existing_field", "Hello World")

	if len(newBuilder.Errors) == 0 {
		t.Fail()
	}

	if newBuilder.Errors[0].Error() != "unable to find index" {
		t.Fail()
	}
}

func TestQueryBuilder_LowerThan(t *testing.T) {
	test := CreateNewQupDynamo(TestConnection{})
	builder, err := test.Query("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	newBuilder := builder.LowerThan("external_id", 123)
	checkVal := newBuilder.Query.KeyConditions["external_id"]

	if *checkVal.AttributeValueList[0].N != "123" {
		t.Fail()
	}

	if *checkVal.ComparisonOperator != "LT" {
		t.Fail()
	}

	if *newBuilder.Query.IndexName != "external_id-group_key" {
		t.Fail()
	}

	newBuilder = builder.LowerThan("non_existing_field", "Hello World")

	if len(newBuilder.Errors) == 0 {
		t.Fail()
	}

	if newBuilder.Errors[0].Error() != "unable to find index" {
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

	builder.Connection = TestConnection{
		MockQueryError: fmt.Errorf("failed due to an unknown curse"),
	}

	_, errs = builder.Count()

	if err == nil || err.Error() != "failed due to an unknown curse" {

	}
}
