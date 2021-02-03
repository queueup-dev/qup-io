package helper_functions_aws_events

// these event wrappers are deprecated, the required functions are available in a map of functions, indexed by reflect.Type.

import (
	"github.com/aws/aws-lambda-go/events"
)

var (
	_ = events.APIGatewayProxyResponse(APIGatewayProxyResponse{})
	_ = events.APIGatewayProxyRequest(APIGatewayProxyRequest{})
	_ = APIGatewayProxyResponse(events.APIGatewayProxyResponse{})
	_ = APIGatewayProxyRequest(events.APIGatewayProxyRequest{})
)

type (
	APIGatewayProxyRequest struct {
		Resource                        string                               `json:"resource"` // The resource path defined in API Gateway
		Path                            string                               `json:"path"`     // The url path for the caller
		HTTPMethod                      string                               `json:"httpMethod"`
		Headers                         map[string]string                    `json:"headers"`
		MultiValueHeaders               map[string][]string                  `json:"multiValueHeaders"`
		QueryStringParameters           map[string]string                    `json:"queryStringParameters"`
		MultiValueQueryStringParameters map[string][]string                  `json:"multiValueQueryStringParameters"`
		PathParameters                  map[string]string                    `json:"pathParameters"`
		StageVariables                  map[string]string                    `json:"stageVariables"`
		RequestContext                  events.APIGatewayProxyRequestContext `json:"requestContext"`
		Body                            string                               `json:"body"`
		IsBase64Encoded                 bool                                 `json:"isBase64Encoded,omitempty"`
	}
	APIGatewayProxyResponse struct {
		StatusCode        int                 `json:"statusCode"`
		Headers           map[string]string   `json:"headers"`
		MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
		Body              string              `json:"body"`
		IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
	}
)

func (q APIGatewayProxyRequest) GetMessage() []byte {
	return []byte(q.Body)
}

func (q *APIGatewayProxyResponse) SetMessage(in []byte) {
	q.Body = string(in)
}
