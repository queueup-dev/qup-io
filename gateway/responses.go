package gateway

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/queueup-dev/qup-io/http"
)

type errorResponse struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}
type Envelope struct {
	Status int16       `json:"status"`
	Data   interface{} `json:"data"`
}

var (
	CORSHeaders = http.Headers{
		"Access-Control-Allow-Origin": "*",
	}
)

func CreateGatewayErrorResponse(message string, err error, status int) (events.APIGatewayProxyResponse, error) {

	response := &errorResponse{
		Type:   "Error",
		Title:  message,
		Status: status,
		Detail: err.Error(),
	}

	body, err := json.Marshal(response)

	if err != nil {
		return CreateGatewayErrorResponse("failed generating error body", err, 500)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
	}, nil
}

func CreateGatewayResponse(body interface{}, headers http.Headers, status int) (events.APIGatewayProxyResponse, error) {

	marshalledBody, err := json.Marshal(Envelope{
		Status: int16(status),
		Data:   body,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    headers,
		Body:       string(marshalledBody),
	}, nil
}
