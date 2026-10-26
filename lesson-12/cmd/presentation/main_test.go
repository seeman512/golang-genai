package main

import (
	"context"
	"testing"
)

func TestHandler(t *testing.T) {
	got, err := handler(context.Background(), Request{Name: "Ada"})
	if err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}
	if got.Message != "Hello, Ada" {
		t.Fatalf("message = %q, want %q", got.Message, "Hello, Ada")
	}
}

func TestHandlerWithEmptyName(t *testing.T) {
	got, err := handler(context.Background(), Request{})
	if err != nil {
		t.Fatalf("handler returned an error: %v", err)
	}
	if got.Message != "Hello, " {
		t.Fatalf("message = %q, want %q", got.Message, "Hello, ")
	}
}
