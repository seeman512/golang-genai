// Package main contains the smallest useful Go Lambda example.
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

// Request is the JSON input accepted by the presentation function.
type Request struct {
	Name string `json:"name"`
}

// Response is the JSON response returned by the presentation function.
type Response struct {
	Message string `json:"message"`
}

// handler contains the business logic. Keeping it separate from main makes it
// possible to test the function without starting the Lambda runtime.
func handler(ctx context.Context, req Request) (Response, error) {
	return Response{
		Message: "Hello, " + req.Name,
	}, nil
}

func main() {
	lambda.Start(handler)
}
