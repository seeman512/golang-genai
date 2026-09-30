package service

import (
	"errors"
	"testing"

	"lesson05/task3_project_structure/internal/repository"
)

func TestAddRejectsBlankTitle(t *testing.T) {
	todoService := New(repository.NewMemoryTodoRepository())

	_, err := todoService.Add(" \t ")
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("Add() error = %v, want ErrEmptyTitle", err)
	}
}

func TestDeleteWrapsNotFound(t *testing.T) {
	todoService := New(repository.NewMemoryTodoRepository())

	err := todoService.Delete(42)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("Delete() error = %v, want wrapped repository.ErrNotFound", err)
	}
}
