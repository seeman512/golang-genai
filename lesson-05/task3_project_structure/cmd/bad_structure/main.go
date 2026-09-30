// Command bad_structure is a small, single-file, in-memory todo list.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type todoItem struct {
	ID        int
	Title     string
	Completed bool
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "todo:", err)
		os.Exit(1)
	}
}

func run() error {
	todos := []todoItem{}
	nextID := 1
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Todo list (items are kept in memory for this run).")
	printHelp()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}

		switch strings.ToLower(parts[0]) {
		case "add":
			if len(parts) < 2 {
				fmt.Println("Usage: add <title>")
				continue
			}

			title := strings.Join(parts[1:], " ")
			todos = append(todos, todoItem{ID: nextID, Title: title})
			fmt.Printf("Added #%d: %s\n", nextID, title)
			nextID++

		case "delete":
			id, ok := readID(parts, "delete")
			if !ok {
				continue
			}

			index := findTodo(todos, id)
			if index == -1 {
				fmt.Printf("No todo item with ID %d.\n", id)
				continue
			}

			fmt.Printf("Deleted #%d: %s\n", todos[index].ID, todos[index].Title)
			todos = append(todos[:index], todos[index+1:]...)

		case "complete":
			id, ok := readID(parts, "complete")
			if !ok {
				continue
			}

			index := findTodo(todos, id)
			if index == -1 {
				fmt.Printf("No todo item with ID %d.\n", id)
				continue
			}

			todos[index].Completed = true
			fmt.Printf("Completed #%d: %s\n", id, todos[index].Title)

		case "list":
			if len(todos) == 0 {
				fmt.Println("No todo items.")
				continue
			}
			for _, item := range todos {
				status := " "
				if item.Completed {
					status = "x"
				}
				fmt.Printf("[%s] #%d %s\n", status, item.ID, item.Title)
			}

		case "help":
			printHelp()

		case "quit", "exit":
			return nil

		default:
			fmt.Printf("Unknown command %q. Type help for available commands.\n", parts[0])
		}
	}

	return scanner.Err()
}

func readID(parts []string, command string) (int, bool) {
	if len(parts) != 2 {
		fmt.Printf("Usage: %s <id>\n", command)
		return 0, false
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil || id < 1 {
		fmt.Println("ID must be a positive integer.")
		return 0, false
	}
	return id, true
}

func findTodo(todos []todoItem, id int) int {
	for i := range todos {
		if todos[i].ID == id {
			return i
		}
	}
	return -1
}

func printHelp() {
	fmt.Println("Commands: add <title>, delete <id>, complete <id>, list, help, quit")
}
