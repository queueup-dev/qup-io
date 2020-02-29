package sentry

import (
	"github.com/getsentry/sentry-go"
)

type Logger struct {
}

func (s Logger) Error(err error, message *string) {
	sentry.CaptureException(err)
}

func (s Logger) Warning(message string, err error) {
	sentry.CaptureMessage(message)
}

func (s Logger) Info(message string) {
	sentry.CaptureMessage(message)
}

func (s Logger) Debug(message string, err error) {
	sentry.CaptureMessage(message)
}
