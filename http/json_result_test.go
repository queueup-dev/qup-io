package http

import (
	"strings"
	"testing"
)

func TestJsonResult_Unmarshal_InvalidJson(t *testing.T) {
	testObject := struct {
		Hello string `json:"hello"`
	}{}

	io := strings.NewReader("{ \"hello: \"world\"}")
	result := &JsonResult{
		output: io,
	}

	err := result.Unmarshal(&testObject)

	if err == nil {
		t.Fail()
	}
}

func TestJsonResult_Unmarshal(t *testing.T) {

	testObject := struct {
		Hello string `json:"hello"`
	}{}

	io := strings.NewReader("{ \"hello\": \"world\"}")
	result := &JsonResult{
		output: io,
	}

	err := result.Unmarshal(&testObject)

	if err != nil {
		t.Error(err)
	}

	if testObject.Hello != "world" {
		t.Failed()
	}
}

func TestJsonResult_ToString(t *testing.T) {
	io := strings.NewReader("{ \"hello\": \"world\"}")
	result := &XmlResult{
		output: io,
	}

	text, _ := result.ToString()

	if *text != "{ \"hello\": \"world\"}" {
		t.Fail()
	}
}
