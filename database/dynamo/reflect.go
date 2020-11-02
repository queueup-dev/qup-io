package dynamo

import (
	"fmt"
	"github.com/queueup-dev/qup-io/reflection"
	"reflect"
	"regexp"
	"strings"
)

type DynamoTableDefinition struct {
	PrimaryKey string
	Fields     []string
	Indices    map[string][]string
}

const (
	primaryIndex   = "key"
	secondaryIndex = "idx"
	rangeKey       = "range"
	dynamoBool     = "B"
	dynamoString   = "S"
	dynamoNumeric  = "N"
)

func tableDefinitionFromStruct(object interface{}) (*DynamoTableDefinition, error) {
	fields, err := reflection.GetTagValues("dynamo", reflect.TypeOf(object))

	if err != nil {
		return nil, err
	}

	definition := &DynamoTableDefinition{
		Indices: map[string][]string{},
	}
	for _, value := range fields {
		parsedValue := strings.Split(value, ",")
		columnName := parsedValue[0]
		definition.Fields = append(definition.Fields, columnName)

		if len(parsedValue) == 1 {
			continue
		}

		typeTag := regexp.MustCompile(`\|`).Split(parsedValue[1], -1)

		switch typeTag[0] {
		case primaryIndex:
			definition.PrimaryKey = columnName
		case secondaryIndex:
			if len(typeTag) < 2 {
				return nil, fmt.Errorf("the %s type should be accompanied with the index name", typeTag[0])
			}

			definition.Indices[typeTag[1]] = append(definition.Indices[typeTag[1]], columnName)
		}
	}

	return definition, nil
}
