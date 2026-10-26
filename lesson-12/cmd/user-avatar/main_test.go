package main

import (
	"context"
	"testing"
)

func TestHandler(t *testing.T) {
	got, err := handler(context.Background(), Request{UserID: 7, FullName: "Ada Lovelace"})
	if err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}
	if got.UserID != 7 || got.Initials != "AL" {
		t.Fatalf("response = %#v, want user 7 with AL initials", got)
	}
	if got.SVG == "" {
		t.Fatal("expected generated SVG")
	}
}

func TestHandlerRejectsMissingName(t *testing.T) {
	if _, err := handler(context.Background(), Request{UserID: 7}); err == nil {
		t.Fatal("handler returned nil error")
	}
}
