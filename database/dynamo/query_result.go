package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type QueryResult struct {
	Result       *dynamodb.QueryOutput
	TargetStruct interface{}
	Decoder      *Decoder
}

func (r QueryResult) First(target interface{}) error {
	for _, item := range r.Result.Items {
		return r.Decoder.UnmarshalMap(item, target)
	}

	return fmt.Errorf("no records found in QueryResult")
}

func (r QueryResult) All(target interface{}) error {
	return r.Decoder.UnmarshalListOfMaps(r.Result.Items, target)
}

func (r QueryResult) Count() *int64 {
	return r.Result.Count
}
