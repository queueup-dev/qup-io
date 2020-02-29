package sentry

import (
	"github.com/getsentry/sentry-go"
	types "github.com/queueup-dev/qup-types"
	"time"
)

type Logger struct {
}

func (s Logger) Error(err error, message *string) {
	sentry.CaptureException(err)
}

func (s Logger) Warning(message string) {
	sentry.CaptureMessage(message)
}

func (s Logger) Info(message string) {
	sentry.CaptureMessage(message)
}

func (s Logger) Debug(message string) {
	sentry.CaptureMessage(message)
}

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
