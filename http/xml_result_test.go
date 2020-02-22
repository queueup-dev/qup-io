package http

import (
	"strings"
	"testing"
)

func TestXmlResult_Unmarshal(t *testing.T) {

	testObject := struct {
		Hello string `xml:"hello"`
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
