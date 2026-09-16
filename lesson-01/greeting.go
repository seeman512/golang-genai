package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const helpMessage = `Usage: greeting [--h] <username>

Print a greeting for a username.

Options:
  --h    show this help message
`

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, output, errorOutput io.Writer) error {
	flags := flag.NewFlagSet("greeting", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	help := flags.Bool("h", false, "show this help message")

	if err := flags.Parse(args); err != nil {
		if helpErr := printHelp(errorOutput); helpErr != nil {
			return fmt.Errorf("print help: %w", helpErr)
		}
		return fmt.Errorf("invalid arguments: %w", err)
	}

	if *help {
		return printHelp(output)
	}

	if flags.NArg() != 1 || strings.TrimSpace(flags.Arg(0)) == "" {
		if helpErr := printHelp(errorOutput); helpErr != nil {
			return fmt.Errorf("print help: %w", helpErr)
		}
		return fmt.Errorf("exactly one non-empty username argument is required")
	}

	if _, err := fmt.Fprintf(output, "Hello, %s!\n", flags.Arg(0)); err != nil {
		return fmt.Errorf("print greeting: %w", err)
	}
	return nil
}

func printHelp(output io.Writer) error {
	_, err := fmt.Fprint(output, helpMessage)
	return err
}
