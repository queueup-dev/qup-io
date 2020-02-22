package http

import (
	"strings"
	"testing"
)

func TestXmlResult_Unmarshal(t *testing.T) {

	testObject := struct {
		Hello string `xml:"hello"`
	}{}

	io := strings.NewReader("<testObject><hello>world</hello></testObject>")
	result := &XmlResult{
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

func TestXmlResult_ToString(t *testing.T) {
	io := strings.NewReader("<testObject><hello>world</hello></testObject>")
	result := &XmlResult{
		output: io,
	}

	text, _ := result.ToString()

	if *text != "<testObject><hello>world</hello></testObject>" {
		t.Fail()
	}
}
