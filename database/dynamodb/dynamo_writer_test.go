package dynamodb

import (
	"encoding/base64"
	"testing"
)

type Testing struct {
	Hello string `json:"hello"`
}

func TestDynamoWriter_Marshal(t *testing.T) {
	writer := DynamoWriter{input: &Testing{Hello: "world"}}

	marshalledContent, err := writer.Bytes()

	if err != nil {
		t.Fail()
	}

	encodedValue := base64.StdEncoding.EncodeToString(marshalledContent)

	if encodedValue != "H/+BAwEBB1Rlc3RpbmcB/4IAAQEBBUhlbGxvAQwAAAAK/4IBBXdvcmxkAA==" {
		t.Fail()
	}
}
