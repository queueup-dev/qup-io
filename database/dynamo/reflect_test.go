package dynamo

import (
	"testing"
)

type ExampleReflect struct {
	Id         string `dynamo:"id,key"`
	ExternalId string `dynamo:"external_id,unique,idx|external_id-group_key"`
	GroupKey   string `dynamo:"group_key,unique,idx|external_id-group_key"`
	OtherField string `dynamo:"other_field,idx|other_field-index"`
}

func TestConditionalExpression(t *testing.T) {
	def, _ := tableDefinitionFromStruct(&ExampleReflect{})

	expression := *conditionExpression(def)

	if expression != "attribute_not_exists(id) AND attribute_not_exists(external_id) AND attribute_not_exists(group_key)" {
		t.Fail()
	}
}
