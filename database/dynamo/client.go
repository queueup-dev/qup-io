package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(s interface{}) error
}

var _ Connection = &dynamodb.DynamoDB{}

type Connection interface {
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
	TransactWriteItems(input *dynamodb.TransactWriteItemsInput) (*dynamodb.TransactWriteItemsOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

type QupDynamo struct {
	Connection Connection
	Validator  Validator
	Decoder    Decoder
	Encoder    Encoder
}

/**
 * Retrieves a single record based on the primaryIndex and loads it in the supplied record.
 */
func (q QupDynamo) Retrieve(table string, key interface{}, record interface{}) error {
	attribute, err := dynamodbattribute.Marshal(key)

	if err != nil {
		return err
	}

	input := &dynamodb.GetItemInput{
		TableName: &table,
		Key: map[string]*dynamodb.AttributeValue{
			primaryIndex: attribute,
		},
	}

	result, err := q.Connection.GetItem(input)

	if err != nil {
		return err
	}

	return q.Decoder.UnmarshalMap(result.Item, &record)
}

/**
 * Starts a (write) Transaction
 * - Save method adds a object to save
 * - Delete method removes an object
 * - Commit commits the changes made in a single transaction
 */
func (q QupDynamo) Transaction(table string, object interface{}) (*TransactionWriter, error) {

	definition, err := tableDefinitionFromStruct(object)

	if err != nil {
		return nil, err
	}

	return &TransactionWriter{
		Connection:      q.Connection,
		TableName:       table,
		TableDefinition: *definition,
		TransactionQuery: &dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{},
		},
		Encoder: &q.Encoder,
	}, nil
}

/**
 * Saves an item to DynamoDb
 */
func (q QupDynamo) Save(table string, record interface{}) error {

	err := q.Validator.Struct(record)

	if err != nil {
		return err
	}

	transaction, err := q.Transaction(table, record)

	if err != nil {
		return err
	}

	transaction.Save(record)
	transaction.Commit()

	return nil
}

/**
 * Query a table. returns a QueryBuilder
 */
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
		Decoder: &q.Decoder,
	}, nil
}

/**
 * Delete a record from table
 */
func (q QupDynamo) Delete(table string, key interface{}) error {

	attribute, err := q.Encoder.Marshal(key)

	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			primaryIndex: attribute,
		},
		TableName: &table,
	}

	_, err = q.Connection.DeleteItem(input)

	return err
}

/**
 * Scan records from a table into target.
 * target should be a map of structures.
 */
func (q QupDynamo) Scan(table string, target interface{}, limit int64) error {

	input := &dynamodb.ScanInput{
		Limit:     &limit,
		TableName: &table,
	}

	output, err := q.Connection.Scan(input)

	if err != nil {
		return err
	}

	return q.Decoder.UnmarshalListOfMaps(output.Items, target)
}

func CreateNewQupDynamo(db Connection) QupDynamo {

	return QupDynamo{
		Connection: db,
		Validator:  validator.New(),
		Decoder:    CreateDecoder(),
		Encoder:    CreateEncoder(),
	}
}
