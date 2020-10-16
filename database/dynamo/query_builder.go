package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/queueup-dev/qup-io/slices"
)

type QueryBuilder struct {
	Connection      *dynamodb.DynamoDB
	TableDefinition DynamoTableDefinition
	Table           string
	Query           *dynamodb.QueryInput
	Errors          []error
	TargetStruct    interface{}
}

func (q QueryBuilder) Equals(field string, value interface{}) QueryBuilder {
	return q.addCondition(field, value, "EQ")
}

func (q QueryBuilder) GreaterThan(field string, value interface{}) QueryBuilder {
	return q.addCondition(field, value, "GT")
}

func (q QueryBuilder) EqualOrLowerThan(field string, value interface{}) QueryBuilder {
	return q.addCondition(field, value, "LE")
}

func (q QueryBuilder) EqualOrGreaterThan(field string, value interface{}) QueryBuilder {
	return q.addCondition(field, value, "GE")
}

func (q QueryBuilder) LowerThan(field string, value interface{}) QueryBuilder {
	return q.addCondition(field, value, "LT")
}

func (q QueryBuilder) addCondition(field string, value interface{}, operator string) QueryBuilder {
	index, err := q.findIndex(field)

	if err != nil {
		return q.addError(err)
	}

	err = q.addIndex(index)

	if err != nil {
		return q.addError(err)
	}

	condition, err := q.createKeyCondition(value, operator)

	if err != nil {
		return q.addError(fmt.Errorf("unable to marshal the value for dynamodb"))
	}

	q.Query.KeyConditions[field] = condition

	return q
}

func (q QueryBuilder) Execute() *[]error {
	if q.Errors != nil && len(q.Errors) != 0 {
		return &q.Errors
	}

	output, err := q.Connection.Query(q.Query)

	if err != nil {
		return &[]error{err}
	}

	fmt.Print(output)
	return nil
}

func (q QueryBuilder) addIndex(index string) error {

	if q.Query.IndexName != nil && *q.Query.IndexName != index {
		return fmt.Errorf("only one index can be queried")
	}

	q.Query.IndexName = &index

	return nil
}

func (q QueryBuilder) addError(err error) QueryBuilder {
	q.Errors = append(q.Errors,
		err,
	)

	return q
}

func (q QueryBuilder) createKeyCondition(value interface{}, operator string) (*dynamodb.Condition, error) {

	dynamodbValue, err := dynamodbattribute.Marshal(value)

	if err != nil {
		return nil, err
	}

	return &dynamodb.Condition{
		AttributeValueList: []*dynamodb.AttributeValue{
			dynamodbValue,
		},
		ComparisonOperator: &operator,
	}, nil
}

func (q QueryBuilder) findIndex(field string) (string, error) {

	if q.TableDefinition.PrimaryKey == field {
		return primaryKey, nil
	}

	for index, val := range q.TableDefinition.GlobalSearchIndices {
		if slices.HasString(field, val) {
			return index, nil
		}
	}

	return "", fmt.Errorf("")
}
