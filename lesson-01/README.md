# greeting

A simple Go CLI that prints a greeting for a given username.

## Usage

```
greeting [--h] <username>
```

Print a greeting for a username.

Options:
  - `--h` — show the help message

## Examples

```bash
go run greeting.go Alice
# Output: Hello, Alice!
```

```bash
go run greeting.go --h
# Output: Usage: greeting [--h] <username>
#
# Print a greeting for a username.
#
# Options:
#   --h    show this help message
```

## Build

```bash
go build -o greeting greeting.go
```

## Test

```bash
go test ./...
```
