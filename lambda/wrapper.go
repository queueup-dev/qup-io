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

// The following two types are added to introduce naming/type conventions.
type InitialWrapper func(handler interface{}) lambda.Handler
type LambdaWrapper func(handler lambda.Handler) (lambda.Handler, error)

// The the Epsagon wrapped handler function is uncommon, and we wrap it separately.
type EpsagonHandler struct {
	epsagonWrappedHandler func(ctx context.Context, payload json.RawMessage) (interface{}, error)
}

// Compile time check that the EpsagonHandler struct implements the lambda.Handler interface.
var _ lambda.Handler = EpsagonHandler{}

func (w EpsagonHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {

	if w.epsagonWrappedHandler == nil {
		panic(fmt.Errorf("handler function is not set"))
	}

	result, err := w.epsagonWrappedHandler(ctx, payload)

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

	return marshalledResult, nil
}

// Start the lambda using the Epsagon wrapper.
// The handler has to adhere to one of the handler type signatures
// as described in the aws documentation of lambda.Start().
func StartLambda(handler interface{}) func() {
	return func() {
		lambda.StartHandler(wrapEpsagon(handler))
	}
}

func wrapHandler(handler interface{}) lambda.Handler {
	return wrapEpsagon(handler)
}

func wrapEpsagon(handler interface{}) lambda.Handler {
	return EpsagonHandler{
		epsagonWrappedHandler: epsagon.WrapLambdaHandler(
			epsagon.NewTracerConfig(envvar.Must("EPSAGON_APP_ID"), envvar.Must("EPSAGON_TOKEN")),
			handler,
		).(func(ctx context.Context, payload json.RawMessage) (interface{}, error)),
	}
}
