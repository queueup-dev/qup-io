package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-playground/validator/v10"
)

type DynamoValidator interface {
	Struct(s interface{}) error
}

type QupDynamo struct {
	Connection *dynamodb.DynamoDB
	Validator  DynamoValidator
}

/**
 * Retrieves a single record based on the primaryKey and loads it in the supplied record.
 */
func (q QupDynamo) Retrieve(table string, key interface{}, record interface{}) error {
	attribute, err := dynamodbattribute.Marshal(key)

	if err != nil {
		return err
	}

	input := &dynamodb.GetItemInput{
		TableName: &table,
		Key: map[string]*dynamodb.AttributeValue{
			primaryKey: attribute,
		},
	}

	result, err := q.Connection.GetItem(input)

	if err != nil {
		return err
	}

	return dynamodbattribute.UnmarshalMap(result.Item, &record)
}

func (q QupDynamo) Save(table string, record interface{}) error {
	err := q.Validator.Struct(record)

	if err != nil {
		return err
	}

	values, err := dynamodbattribute.MarshalMap(record)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      values,
		TableName: &table,
	}

	_, err = q.Connection.PutItem(input)

	return err
}

func (q QupDynamo) Query(table string, object interface{}) (*QueryBuilder, error) {

	definition, err := tableDefinitionFromStruct(object)

	if err != nil {
		return nil, err
	}

	return &QueryBuilder{
		Connection:      q.Connection,
		TableDefinition: *definition,
		TargetStruct:    object,
		Query: &dynamodb.QueryInput{
			TableName:     &table,
			KeyConditions: map[string]*dynamodb.Condition{},
		},
	}, nil
}

//func (q QupDynamo) List(table string, record interface{}) error {
//
//}

func CreateNewQupDynamo(db *dynamodb.DynamoDB) QupDynamo {

	if db == nil {

	}

	return QupDynamo{
		Connection: db,
		Validator:  validator.New(),
	}
}
