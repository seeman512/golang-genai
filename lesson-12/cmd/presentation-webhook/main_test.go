package main

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerProcessesWebhook(t *testing.T) {
	response, err := handler(context.Background(), events.APIGatewayProxyRequest{
		Body: `{"event":"book.created","data":{"id":42}}`,
	})
	if err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", response.StatusCode)
	}
	if response.Body != `{"status":"processed"}` {
		t.Fatalf("body = %q, want processed response", response.Body)
	}
}

func TestHandlerRejectsInvalidJSON(t *testing.T) {
	response, err := handler(context.Background(), events.APIGatewayProxyRequest{Body: "{"})
	if err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}
	if response.StatusCode != 400 {
		t.Fatalf("status = %d, want 400", response.StatusCode)
	}
	if response.Body != `{"error":"invalid json"}` {
		t.Fatalf("body = %q, want invalid-json response", response.Body)
	}
}

func TestHandlerRejectsMissingEvent(t *testing.T) {
	response, err := handler(context.Background(), events.APIGatewayProxyRequest{
		Body: `{"data":{}}`,
	})
	if err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}
	if response.StatusCode != 400 {
		t.Fatalf("status = %d, want 400", response.StatusCode)
	}
}

func TestHandlerDecodesBase64Body(t *testing.T) {
	body := base64.StdEncoding.EncodeToString([]byte(`{"event":"book.created","data":{}}`))
	response, err := handler(context.Background(), events.APIGatewayProxyRequest{
		Body:            body,
		IsBase64Encoded: true,
	})
	if err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", response.StatusCode)
	}
}
