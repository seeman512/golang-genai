// Package handler provides the command-line interface for the todo app.
package handler

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"lesson05/task3_project_structure/internal/service"
)

// Handler translates command-line input into todo service calls.
type Handler struct {
	service *service.Service
	input   io.Reader
	output  io.Writer
}

// New creates a command-line handler using the supplied service and streams.
func New(todoService *service.Service, input io.Reader, output io.Writer) *Handler {
	return &Handler{service: todoService, input: input, output: output}
}

// Run reads and handles commands until the user exits or input ends.
func (h *Handler) Run() error {
	scanner := bufio.NewScanner(h.input)
	if err := h.printf("Todo list (items are kept in memory for this run).\n"); err != nil {
		return err
	}
	if err := h.printHelp(); err != nil {
		return err
	}

	for {
		if err := h.printf("> "); err != nil {
			return err
		}
		if !scanner.Scan() {
			break
		}

		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}

		switch strings.ToLower(fields[0]) {
		case "add":
			if len(fields) < 2 {
				if err := h.printf("Usage: add <title>\n"); err != nil {
					return err
				}
				continue
			}

			item, err := h.service.Add(strings.Join(fields[1:], " "))
			if err != nil {
				if writeErr := h.printf("Error: %v\n", err); writeErr != nil {
					return writeErr
				}
				continue
			}
			if err := h.printf("Added #%d: %s\n", item.ID, item.Title); err != nil {
				return err
			}

		case "delete":
			id, ok, err := h.readID(fields, "delete")
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			if err := h.service.Delete(id); err != nil {
				if writeErr := h.printf("Error: %v\n", err); writeErr != nil {
					return writeErr
				}
				continue
			}
			if err := h.printf("Deleted #%d\n", id); err != nil {
				return err
			}

		case "complete":
			id, ok, err := h.readID(fields, "complete")
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			if err := h.service.Complete(id); err != nil {
				if writeErr := h.printf("Error: %v\n", err); writeErr != nil {
					return writeErr
				}
				continue
			}
			if err := h.printf("Completed #%d\n", id); err != nil {
				return err
			}

		case "list":
			items := h.service.List()
			if len(items) == 0 {
				if err := h.printf("No todo items.\n"); err != nil {
					return err
				}
				continue
			}
			for _, item := range items {
				status := " "
				if item.Completed {
					status = "x"
				}
				if err := h.printf("[%s] #%d %s\n", status, item.ID, item.Title); err != nil {
					return err
				}
			}

		case "help":
			if err := h.printHelp(); err != nil {
				return err
			}

		case "quit", "exit":
			return nil

		default:
			if err := h.printf("Unknown command %q. Type help for available commands.\n", fields[0]); err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read command: %w", err)
	}
	return nil
}

func (h *Handler) readID(fields []string, command string) (int, bool, error) {
	if len(fields) != 2 {
		if err := h.printf("Usage: %s <id>\n", command); err != nil {
			return 0, false, err
		}
		return 0, false, nil
	}

	id, err := strconv.Atoi(fields[1])
	if err != nil || id < 1 {
		if writeErr := h.printf("ID must be a positive integer.\n"); writeErr != nil {
			return 0, false, writeErr
		}
		return 0, false, nil
	}
	return id, true, nil
}

func (h *Handler) printHelp() error {
	return h.printf("Commands: add <title>, delete <id>, complete <id>, list, help, quit\n")
}

func (h *Handler) printf(format string, args ...any) error {
	if _, err := fmt.Fprintf(h.output, format, args...); err != nil {
		return fmt.Errorf("write command output: %w", err)
	}
	return nil
}
