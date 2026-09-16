package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunGreetsUsername(t *testing.T) {
	var output, errorOutput bytes.Buffer

	if err := run([]string{"Alice"}, &output, &errorOutput); err != nil {
		t.Fatalf("run returned an error: %v", err)
	}
	if got, want := output.String(), "Hello, Alice!\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
	if errorOutput.Len() != 0 {
		t.Fatalf("unexpected error output: %q", errorOutput.String())
	}
}

func TestRunHelp(t *testing.T) {
	var output, errorOutput bytes.Buffer

	if err := run([]string{"--h"}, &output, &errorOutput); err != nil {
		t.Fatalf("run returned an error: %v", err)
	}
	if !strings.Contains(output.String(), "Usage: greeting [--h] <username>") {
		t.Fatalf("help output does not contain usage: %q", output.String())
	}
	if errorOutput.Len() != 0 {
		t.Fatalf("unexpected error output: %q", errorOutput.String())
	}
}

func TestRunRequiresOneUsername(t *testing.T) {
	tests := [][]string{
		{},
		{"Alice", "Bob"},
		{"   "},
	}

	for _, args := range tests {
		var output, errorOutput bytes.Buffer
		if err := run(args, &output, &errorOutput); err == nil {
			t.Fatalf("run(%q) returned nil error", args)
		}
		if !strings.Contains(errorOutput.String(), "Usage:") {
			t.Fatalf("run(%q) did not print help: %q", args, errorOutput.String())
		}
	}
}
