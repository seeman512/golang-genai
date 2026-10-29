package server_test

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "mystic-square-grpc/gen/mysticsquare"
	"mystic-square-grpc/internal/auth"
	"mystic-square-grpc/internal/server"
	"mystic-square-grpc/internal/store"
)

var authTestSecret = []byte("lesson-13-test-secret")

func authClient(t *testing.T) pb.MysticSquareServiceClient {
	t.Helper()
	listener := bufconn.Listen(1 << 20)
	authenticator := auth.NewAuthenticator(authTestSecret, auth.DefaultPolicy())
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(authenticator.Unary()), grpc.ChainStreamInterceptor(authenticator.Stream()))
	pb.RegisterMysticSquareServiceServer(grpcServer, server.New(store.New()))
	go func() { _ = grpcServer.Serve(listener) }()
	conn, err := grpc.NewClient("passthrough:///bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close(); grpcServer.Stop(); _ = listener.Close() })
	return pb.NewMysticSquareServiceClient(conn)
}

func authenticatedContext(t *testing.T, role string, ttl time.Duration) context.Context {
	t.Helper()
	token, err := auth.Sign(authTestSecret, auth.Claims{Subject: "test-user", Role: role, ExpiresAt: time.Now().Add(ttl).Unix()})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

func assertCode(t *testing.T, err error, want codes.Code) {
	t.Helper()
	if got := status.Code(err); got != want {
		t.Fatalf("status code = %s, want %s (err=%v)", got, want, err)
	}
}

func TestAuthRequiresToken(t *testing.T) {
	client := authClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.GetGame(ctx, &pb.GetGameRequest{GameId: 1})
	assertCode(t, err, codes.Unauthenticated)
}

func TestAuthRoles(t *testing.T) {
	client := authClient(t)
	player := authenticatedContext(t, auth.RolePlayer, time.Hour)
	game, err := client.CreateGame(player, &pb.CreateGameRequest{Rows: 2, Columns: 2})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.GetGame(player, &pb.GetGameRequest{GameId: game.GetId()}); err != nil {
		t.Fatal(err)
	}
	_, err = client.Register(player, &pb.RegisterRequest{Email: "player@example.com", FullName: "Player"})
	assertCode(t, err, codes.PermissionDenied)
	moderator := authenticatedContext(t, auth.RoleModerator, time.Hour)
	if _, err := client.Register(moderator, &pb.RegisterRequest{Email: "moderator@example.com", FullName: "Moderator"}); err != nil {
		t.Fatal(err)
	}
}

func TestAuthProtectsStreamingRPCs(t *testing.T) {
	client := authClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := client.WatchRegistrationStatus(ctx, &pb.WatchRegistrationRequest{RegistrationId: 1, MaxEvents: 1})
	if err == nil {
		_, err = stream.Recv()
	}
	assertCode(t, err, codes.Unauthenticated)
	player := authenticatedContext(t, auth.RolePlayer, time.Hour)
	moves, err := client.SubmitGameMoves(player)
	if err != nil {
		t.Fatal(err)
	}
	if err := moves.Send(&pb.GameMove{GameId: 999, Row: 0, Column: 0}); err != nil {
		t.Fatal(err)
	}
	_, err = moves.CloseAndRecv()
	// Authentication succeeded; the request now fails for a domain reason.
	assertCode(t, err, codes.NotFound)
}
