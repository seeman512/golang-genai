.PHONY: proto tools certs server gateway client test bench

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0

proto:
	protoc -I proto -I third_party \
		--go_out=. --go_opt=module=mystic-square-grpc \
		--go-grpc_out=. --go-grpc_opt=module=mystic-square-grpc \
		--grpc-gateway_out=. --grpc-gateway_opt=module=mystic-square-grpc \
		proto/mystic_square.proto

certs:
	go run ./cmd/gencerts

server:
	go run ./cmd/mystic-server

gateway:
	go run ./cmd/gateway

client:
	go run ./cmd/mystic-client

test:
	go test ./... -race -count=1

bench:
	go test ./bench -run '^$$' -bench=. -benchmem -benchtime=2s
