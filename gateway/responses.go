package gateway

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type errorResponse struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func CreateGatewayErrorResponse(message string, err error, status int) events.APIGatewayProxyResponse {

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
	}
}
