package handler_test

import (
	"bytes"
	"strings"
	"testing"

	"lesson05/task3_project_structure/internal/handler"
	"lesson05/task3_project_structure/internal/repository"
	"lesson05/task3_project_structure/internal/service"
)

func TestHandlerTodoCommands(t *testing.T) {
	input := strings.NewReader("add Buy groceries\ncomplete 1\nlist\ndelete 1\nlist\nquit\n")
	var output bytes.Buffer

	todoService := service.New(repository.NewMemoryTodoRepository())
	app := handler.New(todoService, input, &output)
	if err := app.Run(); err != nil {
		t.Fatalf("Run() returned an unexpected error: %v", err)
	}

	got := output.String()
	for _, want := range []string{
		"Added #1: Buy groceries",
		"Completed #1",
		"[x] #1 Buy groceries",
		"Deleted #1",
		"No todo items.",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Run() output does not contain %q:\n%s", want, got)
		}
	}
}
