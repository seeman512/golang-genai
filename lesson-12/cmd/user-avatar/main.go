// Package main exposes the deterministic user-avatar Lambda.
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"lesson12/internal/avatar"
)

// Request is the event accepted by the avatar Lambda.
type Request struct {
	UserID   uint64 `json:"user_id"`
	FullName string `json:"full_name"`
}

// Response contains a self-contained avatar preview. A production deployment
// could store the SVG in object storage and return its URL instead.
type Response struct {
	UserID uint64 `json:"user_id"`
	avatar.Result
}

func handler(ctx context.Context, req Request) (Response, error) {
	result, err := avatar.Generate(req.UserID, req.FullName)
	if err != nil {
		return Response{}, err
	}
	return Response{UserID: req.UserID, Result: result}, nil
}

func main() {
	lambda.Start(handler)
}
