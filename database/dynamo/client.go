package dynamo

import (
	"fmt"
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
 * Retrieves an item without unmarshalling.
 */
func (q QupDynamo) retrieveQuery(table string, key interface{}, record interface{}) (*dynamodb.GetItemOutput, error) {
	tableDef, err := tableDefinitionFromStruct(record)
	if err != nil {
		return nil, err
	}

	attribute, err := dynamodbattribute.Marshal(key)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: &table,
		Key: map[string]*dynamodb.AttributeValue{
			tableDef.PrimaryKey.Field: attribute,
		},
	}
	result, err := q.Connection.GetItem(input)

	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
 * Retrieves a single record based on the primaryIndex and loads it in the supplied record.
 */
func (q QupDynamo) Retrieve(table string, key interface{}, record interface{}) error {

	result, err := q.retrieveQuery(table, key, record)

	if err != nil {
		return err
	}

	if len(result.Item) == 0 {
		return fmt.Errorf("record not found")
	}

	return q.Decoder.UnmarshalMap(result.Item, &record)
}

/**
 * Checks whether or not the item exists.
 */
func (q QupDynamo) Exists(table string, key interface{}, record interface{}) (bool, error) {

	result, err := q.retrieveQuery(table, key, record)

	if err != nil {
		return false, err
	}

	if len(result.Item) > 0 {
		return true, nil
	}

	return false, nil
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
 * Updates an item to DynamoDb
 */
func (q QupDynamo) Update(table string, record interface{}) error {

	err := q.Validator.Struct(record)

	if err != nil {
		return err
	}

	transaction, err := q.Transaction(table, record)

	if err != nil {
		return err
	}

	transaction.Update(record)
	errs := transaction.Commit()

	if errs != nil {
		for _, err := range *errs {
			if isConditionalCheckFailedError(err) {
				return ItemDoesNotExistException{Message: "Item does not exist."}
			}
		}
		return fmt.Errorf("something went wrong while updating the record")
	}

	return nil
}

/**
 * Saves an item to DynamoDb
 *
 */
func (q QupDynamo) Create(table string, record interface{}) error {

	err := q.Validator.Struct(record)

	if err != nil {
		return err
	}

	transaction, err := q.Transaction(table, record)

	if err != nil {
		return err
	}

	transaction.Create(record)
	errs := transaction.Commit()

	if errs != nil {
		for _, err := range *errs {
			if isConditionalCheckFailedError(err) {
				return DuplicateEntryException{Message: "Item already exists."}
			}
		}
		return fmt.Errorf("something went wrong while saving the record")
	}

	return nil
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
	errs := transaction.Commit()

	if errs != nil {
		return fmt.Errorf("something went wrong while saving the record")
	}

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
func (q QupDynamo) Delete(table string, key interface{}, record interface{}) error {

	definition, err := tableDefinitionFromStruct(record)
	if err != nil {
		return err
	}

	attribute, err := q.Encoder.Marshal(key)

	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			definition.PrimaryKey.Field: attribute,
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
