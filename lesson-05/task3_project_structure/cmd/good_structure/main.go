// Command good_structure wires the todo application's layers together.
package main

import (
	"fmt"
	"os"

	"lesson05/task3_project_structure/internal/handler"
	"lesson05/task3_project_structure/internal/repository"
	"lesson05/task3_project_structure/internal/service"
)

func main() {
	todoRepository := repository.NewMemoryTodoRepository()
	todoService := service.New(todoRepository)
	app := handler.New(todoService, os.Stdin, os.Stdout)

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "todo: %v\n", err)
		os.Exit(1)
	}
}
