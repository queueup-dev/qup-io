package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Encoder struct {
	Encoder *dynamodbattribute.Encoder
}

func (e Encoder) Marshal(in interface{}) (*dynamodb.AttributeValue, error) {
	return e.Encoder.Encode(in)
}

func (e Encoder) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	av, err := e.Encoder.Encode(in)
	if err != nil || av == nil || av.M == nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	return av.M, nil
}

func (e Encoder) MarshalList(in interface{}) ([]*dynamodb.AttributeValue, error) {
	av, err := e.Encoder.Encode(in)
	if err != nil || av == nil || av.L == nil {
		return []*dynamodb.AttributeValue{}, err
	}

	return av.L, nil
}

func CreateEncoder() Encoder {
	return Encoder{
		Encoder: dynamodbattribute.NewEncoder(func(encoder *dynamodbattribute.Encoder) {
			encoder.TagKey = "dynamo"
		}),
	}
}

func (e Encoder) Must(attribute *dynamodb.AttributeValue, err error) *dynamodb.AttributeValue {
	if err != nil {
		panic(err)
	}

	return attribute
}
