package helper_functions_aws_events

import (
	"testing"
)

func TestAPIGatewayProxyResponse_SetMessage(t *testing.T) {
	var response APIGatewayProxyResponse

	response.SetMessage([]byte("hello"))

	if response.Body != "hello" {
		t.Fail()
	}
}
