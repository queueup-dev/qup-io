package sentry

import (
	"github.com/getsentry/sentry-go"
	types "github.com/queueup-dev/qup-types"
	"time"
)

func NewSentryLogger(dsn string, environment string) types.Logger {
	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 3

	sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: environment,
		Transport:   sentrySyncTransport,
	})

	return &Logger{}
}
