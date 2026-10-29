# Mystic Square: gRPC, grpc-gateway, auth, and mTLS

Lesson 13 is a standalone Go 1.27.1 example of a protobuf API with native gRPC. It also shows how the same contract can be exposed as REST/JSON with grpc-gateway and protected with two independent security layers:

- JWT bearer-token authentication and role-based authorization;
- mutual TLS (mTLS) between the client/gateway and the gRPC server.

The service demonstrates unary calls, server streaming, client streaming, protobuf status codes, cancellation, and clean `io.EOF` stream completion.

## Project layout

```text
proto/mystic_square.proto       API contract and HTTP annotations
third_party/google/api/         protobuf dependencies for grpc-gateway
gen/mysticsquare/                checked-in generated Go and gateway code
cmd/mystic-server/              native gRPC server
cmd/mystic-client/              native gRPC client
cmd/gateway/                    REST/JSON -> gRPC gateway
cmd/mint-token/                 local JWT generator
cmd/gencerts/                   local mTLS certificate generator
internal/auth/                  HS256 demo JWT and gRPC interceptors
internal/tlsconf/               TLS configurations
internal/middleware/            logging and panic recovery
bench/                          gRPC vs REST vs gateway benchmarks
```

## Quick start

The checked-in generated files mean `protoc` is not required for the normal run.

```bash
go test ./...
go run ./cmd/mystic-server
```

In another terminal:

```bash
go run ./cmd/mystic-client register ada@example.com "Ada Lovelace"
go run ./cmd/mystic-client -max 4 watch 1
```

In a third terminal:

```bash
go run ./cmd/mystic-client confirm 1
go run ./cmd/mystic-client complete 1
```

The server starts in plaintext and without authentication only for local development. Use `-auth` and/or `-tls` for the secured examples below.

## Native gRPC client scenarios

Create a game and submit moves through a client stream:

```bash
go run ./cmd/mystic-client create-game 2 2
go run ./cmd/mystic-client moves 1 1,0 1,1
```

The default board is one legal move away from solved. The first move is rejected and the second wins the game. `watch` is a server stream; `moves` is a client stream.

## grpcurl examples

The server enables gRPC reflection so `grpcurl` can inspect and call it without a local `.proto` file. Install it with one of the following:

```bash
brew install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

With the server running on `localhost:50051`:

```bash
# List services and methods.
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 list mysticsquare.v1.MysticSquareService

# Unary calls.
grpcurl -plaintext -d '{"email":"grpcurl@example.com","fullName":"gRPC User"}' \
  localhost:50051 mysticsquare.v1.MysticSquareService/Register
grpcurl -plaintext -d '{"rows":2,"columns":2}' \
  localhost:50051 mysticsquare.v1.MysticSquareService/CreateGame
grpcurl -plaintext -d '{"gameId":1}' \
  localhost:50051 mysticsquare.v1.MysticSquareService/GetGame

# Server-streaming call; maxEvents=4 makes it finish with a clean EOF.
grpcurl -plaintext -d '{"registrationId":1,"maxEvents":4}' \
  localhost:50051 mysticsquare.v1.MysticSquareService/WatchRegistrationStatus

# Client-streaming call: one JSON message per line, then stdin closes.
printf '%s\n' \
  '{"gameId":1,"row":1,"column":0,"sequence":1}' \
  '{"gameId":1,"row":1,"column":1,"sequence":2}' | \
  grpcurl -plaintext -d @ localhost:50051 \
  mysticsquare.v1.MysticSquareService/SubmitGameMoves
```

`grpcurl` uses protobuf JSON names (`fullName`, `gameId`). With PowerShell, use `grpcurl.exe` and escape the JSON quotes as required by the shell.

## REST with grpc-gateway

The HTTP routes are declared beside each RPC in `proto/mystic_square.proto`. Start the native server and then the gateway in separate terminals:

```bash
go run ./cmd/mystic-server
go run ./cmd/gateway
```

The gateway listens on `:8080` and translates JSON to gRPC:

```bash
# Register, create a game, and retrieve it.
curl -i -X POST localhost:8080/v1/registrations \
  -H 'Content-Type: application/json' \
  -d '{"email":"rest@example.com","full_name":"REST User"}'
curl -i -X POST localhost:8080/v1/games \
  -H 'Content-Type: application/json' \
  -d '{"rows":2,"columns":2}'
curl -i localhost:8080/v1/games/1

# Registration transitions.
curl -i -X PATCH localhost:8080/v1/registrations/1/confirm
curl -i -X PATCH localhost:8080/v1/registrations/1/complete

# A server stream is returned as newline-delimited JSON.
curl -N 'localhost:8080/v1/registrations/1/watch?max_events=4'
```

The gateway exposes unary and server-streaming methods. Client-streaming `SubmitGameMoves` has no HTTP mapping and remains native gRPC. The gateway forwards the incoming `Authorization` header to gRPC metadata, so the same token works for REST and gRPC.

For an existing HTTP application, the generated gateway is an `http.Handler`; it can be mounted under a router, for example `router.Handle("/v1/", mux)`.

## Token authentication

The demo JWT is deliberately minimal (HS256, standard library only). It is useful for understanding interceptors, not for production identity management. The default secret is local-demo-only; set `JWT_SECRET` in real deployments.

Start the secured server:

```bash
go run ./cmd/mystic-server -auth
```

Create tokens in another terminal:

```bash
export PLAYER=$(go run ./cmd/mint-token -sub ada -role player)
export MODERATOR=$(go run ./cmd/mint-token -sub grace -role moderator)
```

Use them with the Go client:

```bash
go run ./cmd/mystic-client -token "$PLAYER" get-game 1
go run ./cmd/mystic-client -token "$PLAYER" create-game 2 2
go run ./cmd/mystic-client -token "$PLAYER" register reader@example.com Reader
# error: PermissionDenied (players cannot manage registrations)
go run ./cmd/mystic-client -token "$MODERATOR" register moderator@example.com Moderator
```

Use the token with grpcurl (including reflection discovery when the server runs with `-auth`):

```bash
grpcurl -plaintext -H "authorization: Bearer $PLAYER" \
  list
grpcurl -plaintext -H "authorization: Bearer $PLAYER" \
  -d '{"gameId":1}' localhost:50051 mysticsquare.v1.MysticSquareService/GetGame
```

Use it through the gateway:

```bash
curl -i -H "Authorization: Bearer $PLAYER" localhost:8080/v1/games/1
curl -i -X POST -H "Authorization: Bearer $PLAYER" \
  -H 'Content-Type: application/json' localhost:8080/v1/registrations \
  -d '{"email":"denied@example.com","full_name":"Denied"}'
# HTTP 403: the player is authenticated but not authorized
```

Both unary and stream interceptors are installed. This is important: a unary interceptor alone would leave `WatchRegistrationStatus` and `SubmitGameMoves` unprotected. Roles are configured in `internal/auth/interceptor.go`.

On PowerShell, use `$env:PLAYER = go run ./cmd/mint-token -role player` and `$env:PLAYER` in place of `$PLAYER`.

## mTLS

Generate a local CA, server certificate, trusted client certificate, and an untrusted client certificate:

```bash
go run ./cmd/gencerts
```

Run the gRPC server with mutual TLS and token auth:

```bash
go run ./cmd/mystic-server -tls -auth
```

The native client presents the trusted client certificate and verifies the server certificate:

```bash
go run ./cmd/mystic-client -tls -token "$PLAYER" get-game 1
```

The gateway can also use mTLS for its upstream gRPC connection:

```bash
go run ./cmd/gateway -tls
curl -H "Authorization: Bearer $PLAYER" localhost:8080/v1/games/1
```

The gateway's HTTP listener is local plaintext in this demo; `-tls` secures the gateway-to-gRPC hop. For `grpcurl`, provide the CA and client keypair:

```bash
grpcurl -cacert certs/ca.pem \
  -cert certs/client.pem -key certs/client-key.pem \
  -H "authorization: Bearer $PLAYER" \
  -d '{"gameId":1}' localhost:50051 mysticsquare.v1.MysticSquareService/GetGame
```

Useful negative checks:

```bash
# No client certificate: TLS handshake fails.
go run ./cmd/mystic-client -tls -client-cert none -token "$PLAYER" get-game 1
# Certificate from an untrusted CA: TLS handshake fails.
go run ./cmd/mystic-client -tls -client-cert rogue-client -token "$PLAYER" get-game 1
# Trusted mTLS but no token: TLS succeeds, RPC returns Unauthenticated.
go run ./cmd/mystic-client -tls get-game 1
```

mTLS authenticates the calling process/service; the JWT authenticates and authorizes the user. The two checks are intentionally independent. `certs/` is ignored by Git and is for local demos only.

## Benchmark

The benchmark compares a native gRPC request, hand-written REST/JSON, and REST through grpc-gateway. It uses real loopback TCP connections, reuses connections, and drains/closes HTTP response bodies.

```bash
go test ./bench -run TestPayloadSize -v
go test ./bench -run '^$' -bench=. -benchmem -benchtime=2s
```

Results depend on the machine and Go version. The benchmark is a measurement template, not a claim that one transport always wins. It also includes protobuf-vs-JSON serialization benchmarks.

## Tests and regeneration

```bash
gofmt -w $(find . -name '*.go' -type f)
go vet ./...
go test ./...
```

To regenerate the checked-in protobuf and gateway code, install `protoc` and run:

```bash
make tools
make proto
```

The command requires the `google/api` files under `third_party/`. Do not edit generated files manually. `go mod tidy` is needed after dependency changes.
