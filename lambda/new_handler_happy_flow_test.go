package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/queueup-dev/qup-io/lambda/helper_functions_aws_events"
	"testing"
)

type RequestTestStruct struct {
	Foo int
}

type ResponseTestStruct struct {
	Bar int
}

func (r RequestTestStruct) GetGatewayResponse() (helper_functions_aws_events.APIGatewayProxyResponse, ResponseTestStruct, error) {
	return helper_functions_aws_events.APIGatewayProxyResponse{StatusCode: r.Foo}, ResponseTestStruct{Bar: r.Foo * 2}, nil
}

func TestCompileHandler(t *testing.T) {
	gatewayPayload := []byte("{\"HTTPMethod\":\"awesome-header\",\"Body\":\"{\\\"Foo\\\":123}\"}")
	want := `{"Bar":246}`

	type testRecord struct {
		name string
		args interface{}
	}

	tests := []testRecord{
		{
			name: "no pointers",
			args: func(ctx context.Context, event events.APIGatewayProxyRequest, message RequestTestStruct) (helper_functions_aws_events.APIGatewayProxyResponse, ResponseTestStruct, error) {
				return message.GetGatewayResponse()
			},
		},
		{
			name: "few pointers",
			args: func(ctx context.Context, event helper_functions_aws_events.APIGatewayProxyRequest, message RequestTestStruct) (*events.APIGatewayProxyResponse, ResponseTestStruct, error) {
				x, y, z := message.GetGatewayResponse()
				a := events.APIGatewayProxyResponse(x)

				return &a, y, z
			},
		},
		{
			name: "more pointers",
			args: func(ctx context.Context, event *events.APIGatewayProxyRequest, message *RequestTestStruct) (events.APIGatewayProxyResponse, *ResponseTestStruct, error) {
				x, y, z := message.GetGatewayResponse()
				a := events.APIGatewayProxyResponse(x)
				return a, &y, z
			},
		},
		{
			name: "return has no inner response message",
			args: func(ctx context.Context, event *helper_functions_aws_events.APIGatewayProxyRequest, message *RequestTestStruct) (*helper_functions_aws_events.APIGatewayProxyResponse, error) {
				x, y, z := message.GetGatewayResponse()
				innerPayload, _ := json.Marshal(y)
				x.SetMessage(innerPayload)
				return &x, z
			},
		},
		{
			name: "input has no inner request message",
			args: func(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, *ResponseTestStruct, error) {
				var message RequestTestStruct
				json.Unmarshal([]byte(event.Body), &message)
				x, y, z := message.GetGatewayResponse()
				a := events.APIGatewayProxyResponse(x)
				return &a, &y, z
			},
		},
		{
			name: "no inner request or response messages",
			args: func(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
				var message RequestTestStruct
				json.Unmarshal([]byte(event.Body), &message)
				x, y, z := message.GetGatewayResponse()
				innerPayload, _ := json.Marshal(y)
				x.SetMessage(innerPayload)
				a := events.APIGatewayProxyResponse(x)
				return &a, z
			},
		},
		{
			name: "no request event at all",
			args: func() (*helper_functions_aws_events.APIGatewayProxyResponse, error) {
				var event helper_functions_aws_events.APIGatewayProxyRequest
				json.Unmarshal(gatewayPayload, &event)

				var message RequestTestStruct
				json.Unmarshal([]byte(event.Body), &message)
				x, y, z := message.GetGatewayResponse()
				innerPayload, _ := json.Marshal(y)
				x.SetMessage(innerPayload)
				return &x, z
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHandler(tt.args)
			callGatewayHandler(got, gatewayPayload, want, t)
		})
	}
}

func callGatewayHandler(handler lambda.Handler, payload []byte, want string, t *testing.T) {
	resPayload, err := handler.Invoke(context.Background(), payload)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	res := helper_functions_aws_events.APIGatewayProxyResponse{}
	err = json.Unmarshal(resPayload, &res)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if res.Body != want {
		fmt.Printf("inner payload is not correct, have:\n %v \n", res.Body)
		t.Fail()
	}
}
