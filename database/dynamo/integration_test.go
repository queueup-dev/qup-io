package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"testing"
)

var (
	qupDynamo QupDynamo
	tableName = "Persons"
)

type PersonTest struct {
	Id        string `dynamo:"uid,key"`
	FirstName string `dynamo:"first_name"`
	LastName  string `dynamo:"last_name"`
	City      string `dynamo:"city"`
	Age       int32  `dynamo:"age"`
}

func DeleteTable(connection *dynamodb.DynamoDB) {
	input := &dynamodb.DeleteTableInput{TableName: &tableName}

	_, err := connection.DeleteTable(input)

	if err != nil {
		fmt.Println("Error cleaning up table " + tableName)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func CreateTable(connection *dynamodb.DynamoDB) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("uid"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("uid"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := connection.CreateTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Created the table", tableName)
}

func Setup() {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("eu-west-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return
	}
	dbSvc := dynamodb.New(sess)

	result, err := dbSvc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Tables:")
	for _, table := range result.TableNames {
		if *table == tableName {
			DeleteTable(dbSvc)
		}
	}

	CreateTable(dbSvc)

	qupDynamo = CreateNewQupDynamo(dbSvc)
}

func TestScan(t *testing.T) {

	Setup()

	person1 := &PersonTest{
		Id:        uuid.New().String(),
		FirstName: "John",
		LastName:  "Doe",
		City:      "New York",
		Age:       20,
	}

	person2 := &PersonTest{
		Id:        uuid.New().String(),
		FirstName: "Sara",
		LastName:  "Jones",
		City:      "Washington",
		Age:       29,
	}

	person3 := &PersonTest{
		Id:        uuid.New().String(),
		FirstName: "Peter",
		LastName:  "Jansen",
		City:      "Haarlem",
		Age:       42,
	}

	person4 := &PersonTest{
		Id:        uuid.New().String(),
		FirstName: "Jan",
		LastName:  "de Wit",
		City:      "Zoetermeer",
		Age:       67,
	}

	err := qupDynamo.Save(tableName, person1)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	err = qupDynamo.Save(tableName, person2)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	err = qupDynamo.Save(tableName, person3)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	err = qupDynamo.Save(tableName, person4)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	persons := &[]PersonTest{}
	err = qupDynamo.Scan(tableName, persons, 2)

	if len(*persons) != 2 {
		t.Fail()
		t.Log("expected length of 2, length of " + strconv.Itoa(len(*persons)) + " returned")
	}
}

func TestPersonCRUD(t *testing.T) {

	Setup()

	recordId := uuid.New().String()
	person := &PersonTest{
		Id:        recordId,
		FirstName: "John",
		LastName:  "Doe",
		City:      "New York",
		Age:       20,
	}

	err := qupDynamo.Save(tableName, person)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	result := &PersonTest{}
	err = qupDynamo.Retrieve(tableName, recordId, result)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	if result.FirstName != "John" || result.LastName != "Doe" || result.City != "New York" || result.Age != 20 {
		t.Fail()
	}

	expression := "set age = age + :test"
	values := map[string]interface{}{":test": 2}

	err = qupDynamo.Update(tableName, recordId, PersonTest{}, expression, values)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	result2 := &PersonTest{}
	err = qupDynamo.Retrieve(tableName, recordId, result2)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	if result2.Age != 22 {
		t.Log("failed asserting that age " + strconv.Itoa(int(result2.Age)) + " is age 22")
		t.Fail()
	}

	err = qupDynamo.Delete(tableName, recordId, PersonTest{})

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	result3 := &PersonTest{}
	err = qupDynamo.Retrieve(tableName, recordId, result3)

	if err == nil || err.Error() != "record not found" {
		t.Fail()
	}
}
