package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/queueup-dev/qup-io/slices"
	"strings"
)

type QueryBuilder struct {
	Connection      Connection
	TableDefinition DynamoTableDefinition
	Table           string
	Query           *dynamodb.QueryInput
	Errors          []error
	TargetStruct    interface{}
	Decoder         *Decoder
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

func (q QueryBuilder) Limit(limit int64) QueryBuilder {
	q.Query.Limit = &limit

	return q
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

//func (q QueryBuilder) addFilter(field string, value interface{}, operator string) QueryBuilder {
//
//}

func (q QueryBuilder) Select(fields []string) QueryBuilder {
	q.Query.ProjectionExpression = aws.String(strings.Join(fields, ","))

	return q
}

/**
 * Returns the count in an integer, similar to the count on a QueryResult but more performant if you just care about the count
 * The int64 argument is always returned when the list of []error (2nd argument) is empty.
 */
func (q QueryBuilder) Count() (int64, *[]error) {
	if q.hasErrors() {
		return 0, &q.Errors
	}

	q.Query.Select = aws.String("COUNT")

	output, err := q.Connection.Query(q.Query)

	if err != nil {
		return 0, &[]error{err}
	}

	return *output.Count, nil
}

/**
 * Executes the built query and returns the QueryResult.
 * The QueryResult is always returned when the list of []error (2nd argument) is empty.
 */
func (q QueryBuilder) Execute() (*QueryResult, *[]error) {
	if q.hasErrors() {
		return nil, &q.Errors
	}

	output, err := q.Connection.Query(q.Query)

	if err != nil {
		return nil, &[]error{err}
	}

	return &QueryResult{
		Result:       output,
		TargetStruct: q.TargetStruct,
		Decoder:      q.Decoder,
	}, nil
}

func (q QueryBuilder) hasErrors() bool {
	if q.Errors != nil && len(q.Errors) != 0 {
		return true
	}

	return false
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

	for index, val := range q.TableDefinition.GlobalSearchIndices {
		if slices.HasString(field, val) {
			return index, nil
		}
	}

	return "", fmt.Errorf("unable to find index")
}
