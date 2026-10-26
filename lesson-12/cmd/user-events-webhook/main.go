// Package main exposes the user registration webhook.
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/mail"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"lesson12/internal/avatar"
)

const userRegisteredEvent = "user.registered"

// WebhookPayload is the common event envelope. RawMessage lets each event
// version have its own data structure.
type WebhookPayload struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// UserRegistered is the safe subset of registration data needed for a welcome
// email. Passwords and password hashes are deliberately not accepted.
type UserRegistered struct {
	UserID   uint64 `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

// GreetingEmail represents a preview of the email delivery. Stage 1 does not
// contact an SMTP provider or SES; it emulates delivery in the response.
type GreetingEmail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	HTML    string `json:"html"`
}

// WebhookResponse describes the emulated processing result.
type WebhookResponse struct {
	Status string        `json:"status"`
	Event  string        `json:"event"`
	Email  GreetingEmail `json:"email"`
}

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
	if payload.Event != userRegisteredEvent {
		return jsonResponse(400, map[string]string{"error": "unsupported event"}), nil
	}
	if len(payload.Data) == 0 || string(payload.Data) == "null" {
		return jsonResponse(400, map[string]string{"error": "data is required"}), nil
	}

	var event UserRegistered
	if err := decodeStrict(payload.Data, &event); err != nil {
		return jsonResponse(400, map[string]string{"error": "invalid user.registered data"}), nil
	}
	if err := validateUserRegistered(event); err != nil {
		return jsonResponse(400, map[string]string{"error": err.Error()}), nil
	}

	generated, err := avatar.Generate(event.UserID, event.FullName)
	if err != nil {
		return jsonResponse(400, map[string]string{"error": err.Error()}), nil
	}
	greeting := makeGreeting(event, generated.Initials)

	return jsonResponse(200, WebhookResponse{
		Status: "processed",
		Event:  userRegisteredEvent,
		Email:  greeting,
	}), nil
}

func decodeStrict(data []byte, value any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("multiple json values")
		}
		return err
	}
	return nil
}

func validateUserRegistered(event UserRegistered) error {
	if event.UserID == 0 {
		return fmt.Errorf("user_id must be positive")
	}
	if strings.TrimSpace(event.FullName) == "" {
		return fmt.Errorf("full_name is required")
	}
	address, err := mail.ParseAddress(strings.TrimSpace(event.Email))
	if err != nil || address.Address != strings.TrimSpace(event.Email) {
		return fmt.Errorf("email must be valid")
	}
	return nil
}

func makeGreeting(event UserRegistered, initials string) GreetingEmail {
	name := strings.Join(strings.Fields(event.FullName), " ")
	escapedName := html.EscapeString(name)
	escapedInitials := html.EscapeString(initials)
	return GreetingEmail{
		To:      strings.TrimSpace(event.Email),
		Subject: "Welcome to Mystic Square, " + name + "!",
		Text:    fmt.Sprintf("Hello, %s!\n\nWelcome to Mystic Square. Your account is ready, and your avatar (%s) has been prepared.\n\nHave fun solving puzzles!", name, initials),
		HTML:    fmt.Sprintf("<p>Hello, %s!</p><p>Welcome to <strong>Mystic Square</strong>. Your account is ready, and your avatar (<strong>%s</strong>) has been prepared.</p><p>Have fun solving puzzles!</p>", escapedName, escapedInitials),
	}
}

func jsonResponse(status int, value any) events.APIGatewayProxyResponse {
	body, err := json.Marshal(value)
	if err != nil {
		body = []byte(`{"error":"internal server error"}`)
		status = 500
	}
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}
}

func main() {
	lambda.Start(handler)
}
