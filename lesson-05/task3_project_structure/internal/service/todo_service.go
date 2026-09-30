// Package service implements todo application use cases.
package service

import (
	"errors"
	"fmt"
	"strings"

	"lesson05/task3_project_structure/internal/model"
	"lesson05/task3_project_structure/internal/repository"
)

// ErrEmptyTitle indicates that a todo title contains no non-space characters.
var ErrEmptyTitle = errors.New("todo title cannot be empty")

// Service implements todo-list operations using a repository.
type Service struct {
	repository repository.TodoRepository
}

// New creates a todo service backed by repository.
func New(repository repository.TodoRepository) *Service {
	return &Service{repository: repository}
}

// Add creates a todo item after trimming and validating its title.
func (s *Service) Add(title string) (model.Todo, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return model.Todo{}, fmt.Errorf("add todo: %w", ErrEmptyTitle)
	}

	return s.repository.Add(title), nil
}

// Delete removes a todo item by ID.
func (s *Service) Delete(id int) error {
	if err := s.repository.Delete(id); err != nil {
		return fmt.Errorf("delete todo %d: %w", id, err)
	}
	return nil
}

// Complete marks a todo item as completed.
func (s *Service) Complete(id int) error {
	if err := s.repository.Complete(id); err != nil {
		return fmt.Errorf("complete todo %d: %w", id, err)
	}
	return nil
}

// List returns all todo items.
func (s *Service) List() []model.Todo {
	return s.repository.List()
}
