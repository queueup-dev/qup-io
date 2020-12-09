package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/queueup-dev/qup-io/reflection"
	"reflect"
	"regexp"
	"strings"
)

type TableDefinition struct {
	PrimaryKey *FieldDefinition
	Fields     map[string]*FieldDefinition
	Indices    map[string]*IndexDefinition
}

type IndexDefinition struct {
	Fields []*FieldDefinition
}

type FieldDefinition struct {
	Field  string
	Unique bool
}

const (
	primaryIndex   = "key"
	secondaryIndex = "idx"
	uniqueIndex    = "unique"
	rangeKey       = "range"
	dynamoBool     = "B"
	dynamoString   = "S"
	dynamoNumeric  = "N"
)

var (
	// Not thread safe but doesn't have to be.
	definitions = make(map[string]*TableDefinition)
)

func conditionExpression(definition *TableDefinition) *string {
	var conditions []string

	conditions = append(conditions, fmt.Sprintf("attribute_not_exists(%s)", definition.PrimaryKey.Field))

	for _, val := range definition.Fields {
		if val.Unique {
			conditions = append(conditions, fmt.Sprintf("attribute_not_exists(%s)", val.Field))
		}
	}

	return aws.String(strings.Join(conditions, " AND "))
}

func tableDefinitionFromStruct(object interface{}) (*TableDefinition, error) {
	objectType := reflect.TypeOf(object)

	if definitions[objectType.String()] != nil {
		return definitions[objectType.String()], nil
	}

	fields, err := reflection.GetTagValues("dynamo", objectType)

	if err != nil {
		return nil, err
	}

	definition := &TableDefinition{
		Fields:  map[string]*FieldDefinition{},
		Indices: map[string]*IndexDefinition{},
	}

	for _, value := range fields {
		parsedValue := strings.Split(value, ",")
		columnName := parsedValue[0]

		definition.Fields[columnName] = &FieldDefinition{
			Field: columnName,
		}

		if len(parsedValue) == 1 {
			continue
		}

		for _, tag := range parsedValue {
			typeTag := regexp.MustCompile(`\|`).Split(tag, -1)

			switch typeTag[0] {
			case uniqueIndex:
				definition.Fields[columnName].Unique = true
			case primaryIndex:
				definition.PrimaryKey = &FieldDefinition{
					Field: columnName,
				}
			case secondaryIndex:
				if len(typeTag) < 2 {
					return nil, fmt.Errorf("the %s type should be accompanied with the index name", typeTag[0])
				}

				entry, ok := definition.Indices[typeTag[1]]

				if !ok {
					definition.Indices[typeTag[1]] = &IndexDefinition{
						Fields: []*FieldDefinition{
							{
								Field: columnName,
							},
						},
					}
				} else {
					entry.Fields = append(entry.Fields, &FieldDefinition{
						Field: columnName,
					})
				}
			}
		}
	}

	definitions[objectType.String()] = definition

	return definition, nil
}
