package sentry

// Code copied from client.go from the sentry library.

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"reflect"
	"runtime"
	"strings"
)

func eventFromMessage(client *sentry.Client, message string) *sentry.Event {
	event := sentry.NewEvent()
	event.Message = message

	if client.Options().AttachStacktrace {
		event.Threads = []sentry.Thread{{
			Stacktrace: sentry.NewStacktrace(),
			Crashed:    false,
			Current:    true,
		}}
	}

	return event
}

func eventFromException(exception error) *sentry.Event {
	if exception == nil {
		event := sentry.NewEvent()
		event.Message = fmt.Sprintf("Called %s with nil value", callerFunctionName())
		return event
	}

	stacktrace := sentry.ExtractStacktrace(exception)

	if stacktrace == nil {
		stacktrace = sentry.NewStacktrace()
	}

	cause := exception
	// Handle wrapped errors for github.com/pingcap/errors and github.com/pkg/errors
	if ex, ok := exception.(interface{ Cause() error }); ok {
		if c := ex.Cause(); c != nil {
			cause = c
		}
	}

	event := sentry.NewEvent()
	event.Exception = []sentry.Exception{{
		Value:      cause.Error(),
		Type:       reflect.TypeOf(cause).String(),
		Stacktrace: stacktrace,
	}}
	return event
}

func callerFunctionName() string {
	pcs := make([]uintptr, 1)
	runtime.Callers(4, pcs)
	callersFrames := runtime.CallersFrames(pcs)
	callerFrame, _ := callersFrames.Next()
	return baseName(callerFrame.Function)
}

// baseName returns the symbol name without the package or receiver name.
// It replicates https://golang.org/pkg/debug/gosym/#Sym.BaseName, avoiding a
// dependency on debug/gosym.
func baseName(name string) string {
	if i := strings.LastIndex(name, "."); i != -1 {
		return name[i+1:]
	}
	return name
}
