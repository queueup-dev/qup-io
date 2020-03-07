package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/queueup-dev/qup-types"
)

type Logger struct {
	hub *sentry.Hub
}

func (s Logger) Event(content types.DataNode, event *sentry.Event, level sentry.Level) {
	event.Level = level
	content.ApplyToSentryEvent(event)

	s.hub.CaptureEvent(event)
}

func (s Logger) Error(content types.DataNode, err error, message *string) {
	event := eventFromException(err)

	if message != nil {
		event.Message = *message
	}

	s.Event(content, event, sentry.LevelError)
}

func (s Logger) Warning(content types.DataNode, message string, err error) {
	event := eventFromException(err)
	event.Message = message

	s.Event(content, event, sentry.LevelWarning)
}

func (s Logger) Info(content types.DataNode, message string) {
	event := eventFromMessage(s.hub.Client(), message)
	s.Event(content, event, sentry.LevelInfo)
}

func (s Logger) Debug(content types.DataNode, message string, err error) {
	event := eventFromException(err)
	event.Message = message

	s.Event(content, event, sentry.LevelDebug)
}
