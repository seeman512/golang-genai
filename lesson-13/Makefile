.PHONY: proto tools server client test

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

proto:
	protoc -I proto \
		--go_out=. --go_opt=module=mystic-square-grpc \
		--go-grpc_out=. --go-grpc_opt=module=mystic-square-grpc \
		proto/mystic_square.proto

server:
	go run ./cmd/mystic-server

client:
	go run ./cmd/mystic-client

test:
	go test ./... -race -count=1
