package oauth2

import (
	"fmt"
	qupHttp "github.com/queueup-dev/qup-io/http"
	qupTest "github.com/queueup-dev/qup-io/testing"
	"github.com/queueup-dev/qup-io/writer"
	"sync"
	"testing"
)

var (
	logger = qupTest.StdLogger(1)
)

func TestClientCredentials_AccessToken(t *testing.T) {
	var wg sync.WaitGroup
	mockAPI := qupTest.NewMockApi(t, logger, &wg)

	client := NewClientCredentials(
		"http://localhost:8000/oauth/token",
		"clientId-123",
		"clientSecret-123",
		[]string{
			"email",
			"openId",
		},
	)

	mockAPI.Mock().When("/oauth/token", "POST").RespondWith(
		writer.NewJsonWriter(struct {
			AccessToken string `json:"access_token"`
		}{AccessToken: "test123"}),
		qupHttp.Headers{
			"Content-Type": "application/json",
		},
		200,
	)

	mockAPI.Assert().That("/oauth/token", "POST").RequestHeader("Authorization").Eq(
		"Basic Y2xpZW50SWQtMTIzOmNsaWVudFNlY3JldC0xMjM=",
	)

	mockAPI.Assert().That("/oauth/token", "POST").RequestBody().Eq(
		"grant_type=client_credentials&scope=email+openId",
	)

	go mockAPI.Listen("localhost:8000")

	token, err := client.AccessToken()

	if err != nil {
		fmt.Print(err)
		t.Fail()
	}

	if token.AccessToken != "test123" {
		t.Fail()
	}
}
