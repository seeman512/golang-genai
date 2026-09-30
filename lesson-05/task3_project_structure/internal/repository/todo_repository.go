// Package repository provides todo persistence implementations.
package repository

import (
	"errors"

	"lesson05/task3_project_structure/internal/model"
)

// ErrNotFound indicates that a todo item with the requested ID does not exist.
var ErrNotFound = errors.New("todo item not found")

// TodoRepository defines storage operations used by the todo service.
type TodoRepository interface {
	Add(title string) model.Todo
	Delete(id int) error
	Complete(id int) error
	List() []model.Todo
}
