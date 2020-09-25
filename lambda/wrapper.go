package lambda

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/epsagon/epsagon-go/epsagon"
	"github.com/queueup-dev/qup-io/envvar"
)

func StartLambda(handler interface{}) func() {
	return func() {
		lambda.StartHandler(wrapHandler(handler))
	}
}

func wrapHandler(handler interface{}) lambda.Handler {
	return epsagon.WrapLambdaHandler(
		epsagon.NewTracerConfig(envvar.Must("EPSAGON_APP_ID"), envvar.Must("EPSAGON_TOKEN")),
		handler,
	).(lambda.Handler)
}
