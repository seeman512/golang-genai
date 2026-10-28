package server_test

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	pb "mystic-square-grpc/gen/mysticsquare"
	"mystic-square-grpc/internal/server"
	"mystic-square-grpc/internal/store"
)

func start(t *testing.T) pb.MysticSquareServiceClient {
	t.Helper()
	listener := bufconn.Listen(1 << 20)
	grpcServer := grpc.NewServer()
	pb.RegisterMysticSquareServiceServer(grpcServer, server.New(store.New()))
	go func() { _ = grpcServer.Serve(listener) }()
	conn, err := grpc.NewClient("passthrough:///bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close(); grpcServer.Stop() })
	return pb.NewMysticSquareServiceClient(conn)
}

func callContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func TestRegistrationStatusServerStream(t *testing.T) {
	client := start(t)
	ctx := callContext(t)
	reg, err := client.Register(ctx, &pb.RegisterRequest{Email: "ada@example.com", FullName: "Ada"})
	if err != nil {
		t.Fatal(err)
	}
	stream, err := client.WatchRegistrationStatus(ctx, &pb.WatchRegistrationRequest{RegistrationId: reg.GetId(), MaxEvents: 4})
	if err != nil {
		t.Fatal(err)
	}
	want := []pb.RegistrationStatus{pb.RegistrationStatus_UNREGISTERED, pb.RegistrationStatus_PENDING, pb.RegistrationStatus_MAIL_CONFIRMED, pb.RegistrationStatus_REGISTERED}
	for i := 0; i < 2; i++ {
		event, recvErr := stream.Recv()
		if recvErr != nil || event.GetStatus() != want[i] {
			t.Fatalf("history[%d] = %v, %v", i, event, recvErr)
		}
	}
	if _, err := client.ConfirmRegistration(ctx, &pb.RegistrationRequest{RegistrationId: reg.GetId()}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.CompleteRegistration(ctx, &pb.RegistrationRequest{RegistrationId: reg.GetId()}); err != nil {
		t.Fatal(err)
	}
	for i := 2; i < 4; i++ {
		event, recvErr := stream.Recv()
		if recvErr != nil || event.GetStatus() != want[i] {
			t.Fatalf("future[%d] = %v, %v", i, event, recvErr)
		}
	}
	if _, err := stream.Recv(); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestGameMovesClientStream(t *testing.T) {
	client := start(t)
	ctx := callContext(t)
	game, err := client.CreateGame(ctx, &pb.CreateGameRequest{Rows: 2, Columns: 2})
	if err != nil {
		t.Fatal(err)
	}
	stream, err := client.SubmitGameMoves(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := stream.Send(&pb.GameMove{GameId: game.GetId(), Row: 1, Column: 0, Sequence: 1}); err != nil {
		t.Fatal(err)
	}
	if err := stream.Send(&pb.GameMove{GameId: game.GetId(), Row: 1, Column: 1, Sequence: 2}); err != nil {
		t.Fatal(err)
	}
	summary, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
	}
	if summary.GetReceived() != 2 || summary.GetApplied() != 1 || summary.GetRejected() != 1 || !summary.GetIsWin() {
		t.Fatalf("unexpected summary: %v", summary)
	}
}
