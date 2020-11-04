package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            `{"message": "Hello World in go"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
