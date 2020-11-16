package dynamo

import (
	"testing"
)

func TestEncoder_Marshal(t *testing.T) {

	encoder := CreateEncoder()
	intTest := 123

	attr, err := encoder.Marshal(intTest)

	if err != nil {
		t.Fail()
	}

	if *attr.N != "123" {
		t.Fail()
	}

	stringTest := "Hello World!"
	attr, err = encoder.Marshal(stringTest)

	if err != nil {
		t.Fail()
	}

	if *attr.S != "Hello World!" {
		t.Fail()
	}

	boolTest := true
	attr, err = encoder.Marshal(boolTest)

	if err != nil {
		t.Fail()
	}

	if !*attr.BOOL {
		t.Fail()
	}
}

func TestEncoder_MarshalMap(t *testing.T) {
	encoder := CreateEncoder()
	testMap := map[string]interface{}{
		"stringTest": "helloWorld",
		"boolTest":   true,
		"intTest":    232323,
	}

	marshalledMap, err := encoder.MarshalMap(testMap)

	if err != nil {
		t.Fail()
	}

	if *marshalledMap["stringTest"].S != "helloWorld" {
		t.Fail()
	}

	if *marshalledMap["boolTest"].BOOL != true {
		t.Fail()
	}

	if *marshalledMap["intTest"].N != "232323" {
		t.Fail()
	}
}

func TestEncoder_MarshalList(t *testing.T) {
	encoder := CreateEncoder()
	testList := []interface{}{
		"Test123",
		1234567,
		true,
	}

	marshalledList, err := encoder.MarshalList(testList)

	if err != nil {
		t.Fail()
	}

	if *marshalledList[0].S != "Test123" {
		t.Fail()
	}

	if *marshalledList[1].N != "1234567" {
		t.Fail()
	}

	if *marshalledList[2].BOOL != true {
		t.Fail()
	}
}
