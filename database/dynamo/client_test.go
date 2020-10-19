package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"testing"
)

var ()

type ExampleRecord struct {
	Id         string `dynamo:"id,pkey"`
	ExternalId string `dynamo:"external_id,gsi|external_id-group_key"`
	GroupKey   string `dynamo:"group_key,gsi|external_id-group_key"`
}

type DynamoRecordDefinition []string

func TestValues(t *testing.T) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))

	database := dynamodb.New(sess)
	test := CreateNewQupDynamo(database)
	record := ExampleRecord{}

	builder, err := test.Query("PerfectTiming_Timeslots", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	result, _ := builder.Equals("external_id", "5803861a-d475-4726-b9c5-6eb9f395ad6d").Execute()
	result.First(&record)
	fmt.Print(record)
}
