package avatar

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	got, err := Generate(7, "  Ada   Lovelace ")
	if err != nil {
		t.Fatalf("Generate returned an error: %v", err)
	}
	if got.Initials != "AL" {
		t.Fatalf("initials = %q, want %q", got.Initials, "AL")
	}
	if got.Background == "" || got.Foreground == "" {
		t.Fatal("expected avatar colors")
	}
	if !strings.Contains(got.SVG, ">AL</text>") {
		t.Fatalf("SVG does not contain initials: %s", got.SVG)
	}
	if !strings.Contains(got.SVG, "Ada Lovelace avatar") {
		t.Fatalf("SVG does not contain an accessible label: %s", got.SVG)
	}
}

func TestGenerateIsDeterministic(t *testing.T) {
	first, err := Generate(7, "Ada Lovelace")
	if err != nil {
		t.Fatalf("first Generate returned an error: %v", err)
	}
	second, err := Generate(7, "Ada Lovelace")
	if err != nil {
		t.Fatalf("second Generate returned an error: %v", err)
	}
	if first != second {
		t.Fatalf("same input generated different avatars: %#v != %#v", first, second)
	}
}

func TestGenerateInitials(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "Grace Brewster Hopper", want: "GH"},
		{name: "Plato", want: "PL"},
		{name: "Élodie Durand", want: "ÉD"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Generate(1, test.name)
			if err != nil {
				t.Fatalf("Generate returned an error: %v", err)
			}
			if got.Initials != test.want {
				t.Fatalf("initials = %q, want %q", got.Initials, test.want)
			}
		})
	}
}

func TestGenerateRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name   string
		userID uint64
		full   string
	}{
		{name: "missing user id", userID: 0, full: "Ada Lovelace"},
		{name: "missing name", userID: 1, full: "   "},
		{name: "name too long", userID: 1, full: strings.Repeat("a", maxNameRunes+1)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := Generate(test.userID, test.full); err == nil {
				t.Fatal("Generate returned nil error")
			}
		})
	}
}
