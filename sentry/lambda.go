package sentry

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/getsentry/sentry-go"
	types "github.com/queueup-dev/qup-types"
	"time"
)

var (
	proxyEvent events.APIGatewayProxyRequest
)

func NewLambdaSentryLogger(dsn string, environment string, event events.APIGatewayProxyRequest) types.Logger {
	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 3

	sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: environment,
		Transport:   sentrySyncTransport,
		BeforeSend:  BeforeEventContextCallback,
	})

	proxyEvent = event

	return &Logger{}
}

func BeforeEventContextCallback(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	queryString, _ := json.Marshal(proxyEvent.QueryStringParameters)

	event.Request.QueryString = string(queryString)
	event.Request.Data = proxyEvent.Body

	return event
}
