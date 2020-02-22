package http

import (
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

func TestSuccessXmlRequest(t *testing.T) {
	io := strings.NewReader("<Hello>world</Hello>")
	headers := http.Header{}
	headers.Add("Content-Type", "application/xml")
	result, err := Request(
		&MockClient{
			statusCode: 200,
			body:       strings.NewReader("<ResponseObject><status>success</status></ResponseObject>"),
			header:     headers,
		},
		"POST",
		"http://example.org",
		&map[string]string{
			"Content-Type": "application/xml",
		},
		io,
	)

	if err != nil {
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
	io := strings.NewReader("{ \"hello\": \"world\"}")
	_, err := Request(
		&MockClient{
			statusCode: 400,
			body:       nil,
			header: http.Header{
				"Content-Type": {"application/json"},
			},
		},
		"POST",
		"http://example.org",
		&map[string]string{
			"Content-Type": "application/json",
		},
		io,
	)

	if err == nil {
		t.Fail()
	}

	httpError, ok := err.(HttpError)

	if !ok {
		t.Fail()
	}

	if httpError.statusCode != 400 {
		t.Fail()
	}
}

func TestSuccessRequest(t *testing.T) {

	io := strings.NewReader("{ \"hello\": \"world\"}")
	result, err := Request(
		&MockClient{
			statusCode: 200,
			body:       strings.NewReader("{ \"status\": \"success\"}"),
			header: http.Header{
				"Content-Type": {"application/json"},
			},
		},
		"POST",
		"http://example.org",
		&map[string]string{
			"X-Test-Header": "Hi",
		},
		io,
	)

	if err != nil {
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
