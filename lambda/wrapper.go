package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/epsagon/epsagon-go/epsagon"
	"github.com/epsagon/epsagon-go/protocol"
	"github.com/epsagon/epsagon-go/tracer"
	"github.com/queueup-dev/qup-io/envvar"
)

type WrapHandler struct {
	WrappedInvoke func(ctx context.Context, payload []byte) ([]byte, error)
}

func (w WrapHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	return w.WrappedInvoke(ctx, payload)
}

func StartLambda(handler interface{}) func() {
	return func() {
		lambda.StartHandler(wrapHandler(handler))
	}
}

func wrapHandler(handler interface{}) lambda.Handler {
	return WrapHandler{
		WrappedInvoke: wrapEpsagon(handler),
	}
}

func wrapEpsagon(handler interface{}) func(ctx context.Context, payload []byte) ([]byte, error) {
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		result, err := epsagon.WrapLambdaHandler(
			epsagon.NewTracerConfig(envvar.Must("EPSAGON_APP_ID"), envvar.Must("EPSAGON_TOKEN")),
			handler,
		).(func(ctx context.Context, payload json.RawMessage) (interface{}, error))(ctx, payload)

		if err != nil {
			tracer.AddException(&protocol.Exception{
				Type:    "wrapper",
				Message: fmt.Sprintf("Error in wrapper: error in response: %v", err),
				Time:    tracer.GetTimestamp(),
			})

			return nil, err
		}

		marshalledResult, err := json.Marshal(result)

		if err != nil {
			tracer.AddException(&protocol.Exception{
				Type:    "wrapper",
				Message: fmt.Sprintf("Error in wrapper: failed to convert response: %v", err),
				Time:    tracer.GetTimestamp(),
			})

			return nil, err
		}

		return marshalledResult, err
	}
}
