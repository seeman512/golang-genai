// Package bench compares native gRPC, hand-written REST/JSON, and REST via gateway.
// It deliberately uses real loopback TCP connections and reuses connections.
package bench

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	pb "mystic-square-grpc/gen/mysticsquare"
	"mystic-square-grpc/internal/server"
	"mystic-square-grpc/internal/store"
)

type restGame struct {
	ID    uint64  `json:"id"`
	Board []int32 `json:"board"`
	Moves int32   `json:"moves"`
	Win   bool    `json:"is_win"`
	Over  bool    `json:"is_over"`
}

type environment struct {
	grpcClient pb.MysticSquareServiceClient
	restURL    string
	gatewayURL string
	httpClient *http.Client
}

func setup(tb testing.TB) *environment {
	tb.Helper()
	st := store.New()
	game, err := st.CreateGame(2, 2)
	if err != nil {
		tb.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		tb.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMysticSquareServiceServer(grpcServer, server.New(st))
	go func() { _ = grpcServer.Serve(listener) }()
	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		tb.Fatal(err)
	}

	restMux := http.NewServeMux()
	restMux.HandleFunc("GET /games/{id}", func(w http.ResponseWriter, r *http.Request) {
		current, getErr := st.GetGame(game.ID)
		if getErr != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(restGame{ID: uint64(current.ID), Board: board(current.Board), Moves: int32(len(current.Moves)), Win: current.IsWin, Over: current.IsOver})
	})
	restServer := httptest.NewServer(restMux)

	gatewayMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{MarshalOptions: protojson.MarshalOptions{UseProtoNames: true, EmitUnpopulated: true}}))
	if err := pb.RegisterMysticSquareServiceHandler(context.Background(), gatewayMux, conn); err != nil {
		tb.Fatal(err)
	}
	gatewayServer := httptest.NewServer(gatewayMux)
	httpClient := &http.Client{Transport: &http.Transport{MaxIdleConns: 128, MaxIdleConnsPerHost: 128}}
	tb.Cleanup(func() {
		gatewayServer.Close()
		restServer.Close()
		httpClient.CloseIdleConnections()
		_ = conn.Close()
		grpcServer.Stop()
		_ = listener.Close()
	})
	return &environment{grpcClient: pb.NewMysticSquareServiceClient(conn), restURL: restServer.URL + "/games/1", gatewayURL: gatewayServer.URL + "/v1/games/1", httpClient: httpClient}
}

func board(values []int) []int32 {
	out := make([]int32, len(values))
	for i, value := range values {
		out[i] = int32(value)
	}
	return out
}

func (e *environment) getHTTP(tb testing.TB, url string) {
	response, err := e.httpClient.Get(url)
	if err != nil {
		tb.Fatal(err)
	}
	_, _ = io.Copy(io.Discard, response.Body)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK {
		tb.Fatalf("HTTP status = %d", response.StatusCode)
	}
}

func BenchmarkGetGameGRPC(b *testing.B) {
	env := setup(b)
	request := &pb.GetGameRequest{GameId: 1}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := env.grpcClient.GetGame(context.Background(), request); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetGameREST(b *testing.B) {
	env := setup(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env.getHTTP(b, env.restURL)
	}
}

func BenchmarkGetGameRESTViaGateway(b *testing.B) {
	env := setup(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env.getHTTP(b, env.gatewayURL)
	}
}

func BenchmarkSerializeProtobuf(b *testing.B) {
	game := &pb.GameSnapshot{Id: 1, Board: []int32{1, 2, 3, 0}, Moves: 0}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		data, err := proto.Marshal(game)
		if err != nil {
			b.Fatal(err)
		}
		var decoded pb.GameSnapshot
		if err := proto.Unmarshal(data, &decoded); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSerializeJSON(b *testing.B) {
	game := restGame{ID: 1, Board: []int32{1, 2, 3, 0}}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		data, err := json.Marshal(game)
		if err != nil {
			b.Fatal(err)
		}
		var decoded restGame
		if err := json.Unmarshal(data, &decoded); err != nil {
			b.Fatal(err)
		}
	}
}

func TestPayloadSize(t *testing.T) {
	game := &pb.GameSnapshot{Id: 1, Board: []int32{1, 2, 3, 0}}
	protobufBytes, err := proto.Marshal(game)
	if err != nil {
		t.Fatal(err)
	}
	jsonBytes, err := json.Marshal(restGame{ID: 1, Board: []int32{1, 2, 3, 0}})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("protobuf: %d bytes; JSON: %d bytes", len(protobufBytes), len(jsonBytes))
	if len(protobufBytes) >= len(jsonBytes) {
		t.Fatalf("protobuf should be smaller")
	}
}
