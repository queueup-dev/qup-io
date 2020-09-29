package http

import (
	"github.com/queueup-dev/qup-io/writer"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type MockClient struct {
	statusCode int
	body       io.Reader
	header     http.Header
}

func (m MockClient) Do(r *http.Request) (*http.Response, error) {

	response := ioutil.NopCloser(
		m.body,
	)

	return &http.Response{
		Status:     "OK",
		StatusCode: m.statusCode,
		Body:       response,
		Header:     m.header,
	}, nil
}

type ResponseObject struct {
	Status string `json:"status" xml:"status"`
}

type HelloWorld struct {
	Hello string `xml:"Hello" json:"hello"`
}

func TestSuccessXmlRequest(t *testing.T) {
	xmlWriter := writer.NewXmlWriter(&HelloWorld{Hello: "world"})

	headers := http.Header{}
	headers.Add("Content-Type", "application/xml")
	result, _, httpError, err := Request(
		&MockClient{
			statusCode: 200,
			body:       strings.NewReader("<ResponseObject><status>success</status></ResponseObject>"),
			header:     headers,
		},
		"POST",
		"http://example.org",
		&Headers{
			"Content-Type": "application/xml",
		},
		xmlWriter,
	)

	if err != nil || httpError != nil {
		t.Fatal(err)
	}

	object := &ResponseObject{}
	err = result.Unmarshal(object)

	if err != nil {
		t.Fatal(err)
	}

	if object.Status != "success" {
		t.Fail()
	}
}

func TestFailedRequest(t *testing.T) {
	jsonWriter := writer.NewJsonWriter(&HelloWorld{Hello: "world"})
	_, _, httpError, err := Request(
		&MockClient{
			statusCode: 400,
			body:       strings.NewReader("{ \"status\": \"failed\"}"),
			header: http.Header{
				"Content-Type": {"application/json"},
			},
		},
		"POST",
		"http://example.org",
		&Headers{
			"Content-Type": "application/json",
		},
		jsonWriter,
	)

	if err != nil || httpError == nil {
		t.Fail()
	}

	if httpError.StatusCode() != 400 {
		t.Fail()
	}

	message, err := httpError.HttpResponse().ToString()
	if *message != "{ \"status\": \"failed\"}" {
		t.Fail()
	}
}

func TestSuccessRequest(t *testing.T) {

	jsonWriter := writer.NewJsonWriter(&HelloWorld{Hello: "world"})
	result, _, httpError, err := Request(
		&MockClient{
			statusCode: 200,
			body:       strings.NewReader("{ \"status\": \"success\"}"),
			header: http.Header{
				"Content-Type": {"application/json"},
			},
		},
		"POST",
		"http://example.org",
		&Headers{
			"X-Test-Header": "Hi",
		},
		jsonWriter,
	)

	if err != nil || httpError != nil {
		t.Fatal(err)
	}

	object := &ResponseObject{}
	err = result.Unmarshal(object)

	if err != nil {
		t.Fatal(err)
	}

	if object.Status != "success" {
		t.Fail()
	}
}
