package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Decoder struct {
	Decoder *dynamodbattribute.Decoder
}

func (c Decoder) UnmarshalList(l []*dynamodb.AttributeValue, out interface{}) error {
	return c.Decoder.Decode(&dynamodb.AttributeValue{L: l}, out)
}

func (c Decoder) UnmarshalListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	items := make([]*dynamodb.AttributeValue, len(l))
	for i, m := range l {
		items[i] = &dynamodb.AttributeValue{M: m}
	}

	return c.UnmarshalList(items, out)
}

func (c Decoder) UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error {
	return c.Decoder.Decode(&dynamodb.AttributeValue{M: m}, out)
}

func CreateDecoder() Decoder {
	return Decoder{
		Decoder: dynamodbattribute.NewDecoder(func(decoder *dynamodbattribute.Decoder) {
			decoder.TagKey = "dynamo"
		}),
	}
}
