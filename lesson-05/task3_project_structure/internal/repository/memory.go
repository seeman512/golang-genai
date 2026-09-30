package repository

import "lesson05/task3_project_structure/internal/model"

// MemoryTodoRepository stores todo items in memory.
type MemoryTodoRepository struct {
	items  []model.Todo
	nextID int
}

// NewMemoryTodoRepository creates an empty in-memory todo repository.
func NewMemoryTodoRepository() *MemoryTodoRepository {
	return &MemoryTodoRepository{nextID: 1}
}

// Add creates and stores a todo item.
func (r *MemoryTodoRepository) Add(title string) model.Todo {
	if r.nextID < 1 {
		r.nextID = 1
	}

	item := model.Todo{ID: r.nextID, Title: title}
	r.nextID++
	r.items = append(r.items, item)
	return item
}

// Delete removes a todo item by ID.
func (r *MemoryTodoRepository) Delete(id int) error {
	index := r.findIndex(id)
	if index == -1 {
		return ErrNotFound
	}

	r.items = append(r.items[:index], r.items[index+1:]...)
	return nil
}

// Complete marks a todo item as completed.
func (r *MemoryTodoRepository) Complete(id int) error {
	index := r.findIndex(id)
	if index == -1 {
		return ErrNotFound
	}

	r.items[index].Completed = true
	return nil
}

// List returns a copy of all stored todo items.
func (r *MemoryTodoRepository) List() []model.Todo {
	return append([]model.Todo(nil), r.items...)
}

func (r *MemoryTodoRepository) findIndex(id int) int {
	for i := range r.items {
		if r.items[i].ID == id {
			return i
		}
	}
	return -1
}
