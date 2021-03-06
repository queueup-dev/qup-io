package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type TransactionWriter struct {
	Connection       Connection
	TableName        string
	TableDefinition  TableDefinition
	TransactionQuery *dynamodb.TransactWriteItemsInput
	Errors           []error
	Encoder          *Encoder
}

func (b TransactionWriter) Delete(key interface{}) TransactionWriter {

	dynamodbValue, err := b.Encoder.Marshal(key)

	if err != nil {
		return b.addError(err)
	}

	query := dynamodb.TransactWriteItem{
		Delete: &dynamodb.Delete{
			Key: map[string]*dynamodb.AttributeValue{
				b.TableDefinition.PrimaryKey.Field: dynamodbValue,
			},
			TableName: &b.TableName,
		},
	}

	err = b.addInTransaction(query)

	return b.addError(err)
}

func (b TransactionWriter) Save(record interface{}) TransactionWriter {
	values, err := b.Encoder.MarshalMap(record)

	if err != nil {
		return b.addError(err)
	}

	query := dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			Item:      values,
			TableName: &b.TableName,
		},
	}

	err = b.addInTransaction(query)

	return b.addError(err)
}

func (b TransactionWriter) Create(record interface{}) TransactionWriter {

	values, err := b.Encoder.MarshalMap(record)

	if err != nil {
		return b.addError(err)
	}

	query := dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			Item:                values,
			TableName:           &b.TableName,
			ConditionExpression: conditionExpression(&b.TableDefinition),
		},
	}

	err = b.addInTransaction(query)

	return b.addError(err)
}

func (b TransactionWriter) Update(key interface{}, expression string, values map[string]interface{}) TransactionWriter {

	dynamodbValue, err := b.Encoder.Marshal(key)

	if err != nil {
		return b.addError(err)
	}

	mapped, err := b.Encoder.MarshalMap(values)

	if err != nil {
		return b.addError(err)
	}

	expressionAttributeNames := mapExpressionNames(expression)

	if len(expressionAttributeNames) == 0 {
		return b.addError(fmt.Errorf("no expression attribute names found, did you escape it with a \"#\""))
	}

	query := dynamodb.TransactWriteItem{
		Update: &dynamodb.Update{
			ExpressionAttributeNames:  expressionAttributeNames,
			ExpressionAttributeValues: mapped,
			TableName:                 &b.TableName,
			Key: map[string]*dynamodb.AttributeValue{
				b.TableDefinition.PrimaryKey.Field: dynamodbValue,
			},
			UpdateExpression: &expression,
		},
	}

	err = b.addInTransaction(query)

	return b
}

func (b TransactionWriter) SaveExisting(record interface{}) TransactionWriter {

	values, err := b.Encoder.MarshalMap(record)

	if err != nil {
		return b.addError(err)
	}

	query := dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			Item:      values,
			TableName: &b.TableName,
			ExpressionAttributeNames: map[string]*string{
				"#primaryKey": &b.TableDefinition.PrimaryKey.Field,
			},
			ConditionExpression: aws.String("attribute_exists(#primaryKey)"),
		},
	}

	err = b.addInTransaction(query)

	return b.addError(err)
}

func (b TransactionWriter) Commit() *[]error {

	if b.Errors != nil && len(b.Errors) > 0 {
		return &b.Errors
	}

	_, err := b.Connection.TransactWriteItems(b.TransactionQuery)

	if err != nil {
		return &[]error{err}
	}

	return nil
}

func (b TransactionWriter) addInTransaction(item dynamodb.TransactWriteItem) error {
	if len(b.TransactionQuery.TransactItems) > 25 {
		return fmt.Errorf("only up to 25 items are allowed in a single transaction")
	}

	b.TransactionQuery.TransactItems = append(b.TransactionQuery.TransactItems, &item)

	return nil
}

func (b TransactionWriter) addError(err error) TransactionWriter {
	if err != nil {
		b.Errors = append(b.Errors,
			err,
		)
	}

	return b
}
