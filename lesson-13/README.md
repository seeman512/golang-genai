# Mystic Square gRPC / Protobuf Demo

This standalone Lesson 13 project demonstrates protobuf and native gRPC with a small Mystic Square domain. It is intentionally separate from the production HTTP application in `https://github.com/seeman512/mystic-square`.

The service demonstrates:

- unary registration and game calls;
- server streaming for registration status changes;
- client streaming for a sequence of game moves;
- protobuf-generated messages and gRPC status codes;
- stream cancellation and normal `io.EOF` completion.

The registration lifecycle is:

```text
UNREGISTERED -> PENDING -> MAIL_CONFIRMED -> REGISTERED
```

Email delivery is simulated by the `confirm` command. Registration is completed by the `complete` command.

## Run

Go 1.27.1 is required.

```bash
go mod tidy
go run ./cmd/mystic-server
```

In another terminal:

```bash
go run ./cmd/mystic-client register ada@example.com "Ada Lovelace"
# note the returned id
go run ./cmd/mystic-client -max 4 watch 1
```

In a third terminal, advance the registration:

```bash
go run ./cmd/mystic-client confirm 1
go run ./cmd/mystic-client complete 1
```

The watcher receives the replayed initial events and subsequent changes. Without `-max`, stop a watcher with Ctrl+C.

## Client-streaming game moves

Create a game and note its ID:

```bash
go run ./cmd/mystic-client create-game 2 2
```

The default 2×2 game starts as `[1, 0, 3, 2]`. Submit moves as one client stream. The first move below is invalid and the second solves the game:

```bash
go run ./cmd/mystic-client moves 1 1,0 1,1
```

The server applies moves in order and returns one summary after the client half-closes the stream. Rejected moves are included in the summary; an invalid move does not terminate the whole stream.

## Regenerate protobuf code

The checked-in generated files allow the demo to build without protoc. If protoc and its Go plugins are installed:

```bash
make proto
```

Generated files are in `gen/mysticsquare/` and the source contract is `proto/mystic_square.proto`.

## Validation

```bash
gofmt -w $(find . -name '*.go' -type f)
go vet ./...
go test ./...
```
