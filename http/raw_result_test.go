package http

import (
	"strings"
	"testing"
)

func TestRawResult_Unmarshal(t *testing.T) {

	testObject := struct {
		Hello string `xml:"hello"`
	}{}

	io := strings.NewReader("Hello I'm muzzy")
	result := &RawResult{
		output: io,
	}

	err := result.Unmarshal(&testObject)

	if err == nil {
		t.Fail()
	}

	if err.Error() != "unable to unmarshal plain text" {
		t.Fail()
	}

	toString, _ := result.ToString()

	if *toString != "Hello I'm muzzy" {
		t.Fail()
	}
}

func TestRawResult_Casting(t *testing.T) {
	testObject := struct {
		Hello string `json:"hello"`
	}{}

	io := strings.NewReader("{ \"hello\": \"world\"}")
	result := &RawResult{
		output: io,
	}

	jsonResult := JsonResult(*result)
	err := jsonResult.Unmarshal(&testObject)

	if err != nil {
		t.Fatal(err)
	}

	if testObject.Hello != "world" {
		t.Fail()
	}
}
