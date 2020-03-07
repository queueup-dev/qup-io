package sentry

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/getsentry/sentry-go"
	types "github.com/queueup-dev/qup-types"
	"io"
	"strings"
	"time"
)

type ProxyEvent events.APIGatewayProxyRequest

func (e ProxyEvent) ApplyToSentryEvent(event *sentry.Event) {

	if e.QueryStringParameters != nil {
		queryString, _ := json.Marshal(e.QueryStringParameters)
		event.Request.QueryString = string(queryString)
	}

	event.Request.Data = e.Body
	event.Request.Headers = e.Headers
}

func (e ProxyEvent) GetPayload() io.Reader {
	return strings.NewReader(e.Body)
}

func NewLambdaSentryLogger(dsn string, environment string) types.Logger {
	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 3

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: environment,
		Transport:   sentrySyncTransport,
	})

	if err != nil {
		fmt.Println(err)
	}

	return &Logger{
		hub: sentry.CurrentHub(),
	}
}
