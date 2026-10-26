// Package main contains an API Gateway/Lambda proxy webhook example.
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// WebhookPayload is the envelope accepted by the webhook. RawMessage keeps
// the example useful when different event types have different data shapes.
type WebhookPayload struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// handler handles an API Gateway Lambda-proxy request. Client errors are
// returned as HTTP responses with a nil Go error; returning an error here
// would make API Gateway turn an expected 400 into a 502.
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := req.Body
	if req.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return jsonResponse(400, map[string]string{"error": "invalid base64 body"}), nil
		}
		body = string(decoded)
	}

	var payload WebhookPayload
	if strings.TrimSpace(body) == "" || json.Unmarshal([]byte(body), &payload) != nil {
		return jsonResponse(400, map[string]string{"error": "invalid json"}), nil
	}
	if strings.TrimSpace(payload.Event) == "" {
		return jsonResponse(400, map[string]string{"error": "event is required"}), nil
	}
	if len(payload.Data) == 0 || string(payload.Data) == "null" {
		return jsonResponse(400, map[string]string{"error": "data is required"}), nil
	}
	if !json.Valid(payload.Data) {
		return jsonResponse(400, map[string]string{"error": "data must be valid json"}), nil
	}

	// In a real application payload.Event would select the processing path.
	// The lesson keeps that business operation deliberately small.
	return jsonResponse(200, map[string]string{"status": "processed"}), nil
}

func jsonResponse(status int, value any) events.APIGatewayProxyResponse {
	body, err := json.Marshal(value)
	if err != nil {
		// All values used by this example are JSON-marshalable. Keep the fallback
		// here so a future change cannot accidentally return an empty response.
		body = []byte(`{"error":"internal server error"}`)
		status = 500
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}

func main() {
	lambda.Start(handler)
}
