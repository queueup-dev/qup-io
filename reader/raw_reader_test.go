package reader

import (
	"strings"
	"testing"
)

func TestRawResult_Unmarshal(t *testing.T) {

	testObject := struct {
		Hello string `xml:"hello"`
	}{}

	io := strings.NewReader("Hello I'm muzzy")
	result := &RawReader{
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

func TestRawResult_Valid(t *testing.T) {
	io := strings.NewReader("Hello World")
	result := &RawReader{
		output: io,
	}

	if !result.Valid() {
		t.Fail()
	}
}

func TestRawResult_Casting(t *testing.T) {
	testObject := struct {
		Hello string `json:"hello"`
	}{}

	io := strings.NewReader("{ \"hello\": \"world\"}")
	result := &RawReader{
		output: io,
	}

	jsonResult := JsonReader(*result)
	err := jsonResult.Unmarshal(&testObject)

	if err != nil {
		t.Fatal(err)
	}

	if testObject.Hello != "world" {
		t.Fail()
	}
}
