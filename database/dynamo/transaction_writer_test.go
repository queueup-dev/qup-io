package dynamo

import (
	"fmt"
	"testing"
)

func TestTransactionWriter_Delete(t *testing.T) {
	connection := TestConnection{}
	dynamoClient := CreateNewQupDynamo(connection)

	transaction, err := dynamoClient.Transaction("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	newTransaction := transaction.Delete("test123")
	if *newTransaction.TransactionQuery.TransactItems[0].Delete.Key["id"].S != "test123" {
		t.Fail()
	}
}

func TestTransactionWriter_Save(t *testing.T) {
	connection := TestConnection{}
	dynamoClient := CreateNewQupDynamo(connection)

	transaction, err := dynamoClient.Transaction("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	newTransaction := transaction.Save(&ExampleRecord{Id: "test456"})
	if *newTransaction.TransactionQuery.TransactItems[0].Put.Item["id"].S != "test456" {
		t.Fail()
	}
}

func TestTransactionWriter_Commit(t *testing.T) {
	connection := TestConnection{
		MockTransactError: nil,
	}
	dynamoClient := CreateNewQupDynamo(connection)

	transaction, err := dynamoClient.Transaction("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	errs := transaction.Commit()

	if errs != nil {
		t.Fail()
	}
}

func TestTransactionWriter_CommitFail(t *testing.T) {
	connection := TestConnection{
		MockTransactError: fmt.Errorf("transact fail"),
	}
	dynamoClient := CreateNewQupDynamo(connection)

	transaction, err := dynamoClient.Transaction("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	errs := transaction.Commit()

	if errs == nil {
		t.Fail()
	}

	listErrors := *errs

	if len(listErrors) != 1 {
		t.Fail()
	}

	if listErrors[0].Error() != "transact fail" {
		t.Fail()
	}
}

func TestTransactionWriter_Update(t *testing.T) {
	connection := TestConnection{
		MockTransactError: nil,
	}
	dynamoClient := CreateNewQupDynamo(connection)

	transaction, err := dynamoClient.Transaction("mockTable", ExampleRecord{})

	if err != nil {
		t.Fail()
	}

	testValues := map[string]string{
		":test": "1234",
	}

	writer := transaction.Update("12345", "set test = :test", testValues)

	if *writer.TransactionQuery.TransactItems[0].Update.UpdateExpression != "set test = :test" {
		t.Fail()
	}

	if *writer.TransactionQuery.TransactItems[0].Update.ExpressionAttributeValues[":test"].S != "1234" {
		t.Fail()
	}
}
