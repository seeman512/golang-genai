// Package model contains the todo domain model.
package model

// Todo is one item on a todo list.
type Todo struct {
	ID        int
	Title     string
	Completed bool
}
