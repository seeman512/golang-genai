package main

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerEmulatesWelcomeEmail(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		Body: `{"event":"user.registered","data":{"user_id":7,"full_name":"Ada Lovelace","email":"ada@example.com"}}`,
	}
	response, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", response.StatusCode)
	}
	if !containsAll(response.Body, `"status":"processed"`, `"to":"ada@example.com"`, "Welcome to Mystic Square") {
		t.Fatalf("response does not contain the greeting email: %s", response.Body)
	}
}

func TestHandlerAcceptsBase64Body(t *testing.T) {
	body := base64.StdEncoding.EncodeToString([]byte(`{"event":"user.registered","data":{"user_id":7,"full_name":"Ada Lovelace","email":"ada@example.com"}}`))
	response, err := handler(context.Background(), events.APIGatewayProxyRequest{Body: body, IsBase64Encoded: true})
	if err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", response.StatusCode)
	}
}

func TestHandlerRejectsUnsupportedEventAndSensitiveData(t *testing.T) {
	tests := []string{
		`{"event":"game.completed","data":{}}`,
		`{"event":"user.registered","data":{"user_id":7,"full_name":"Ada","email":"ada@example.com","password":"secret"}}`,
	}
	for _, body := range tests {
		response, err := handler(context.Background(), events.APIGatewayProxyRequest{Body: body})
		if err != nil {
			t.Fatalf("handler returned an error: %v", err)
		}
		if response.StatusCode != 400 {
			t.Fatalf("status = %d, want 400 for %s", response.StatusCode, body)
		}
	}
}

func TestHandlerRejectsInvalidRegistration(t *testing.T) {
	response, err := handler(context.Background(), events.APIGatewayProxyRequest{
		Body: `{"event":"user.registered","data":{"user_id":7,"full_name":"Ada","email":"not-an-email"}}`,
	})
	if err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}
	if response.StatusCode != 400 {
		t.Fatalf("status = %d, want 400", response.StatusCode)
	}
}

func containsAll(value string, needles ...string) bool {
	for _, needle := range needles {
		if !contains(value, needle) {
			return false
		}
	}
	return true
}

func contains(value, needle string) bool {
	for i := 0; i+len(needle) <= len(value); i++ {
		if value[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
