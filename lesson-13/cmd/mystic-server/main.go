// mystic-server runs the Lesson 13 gRPC service.
//
//	go run ./cmd/mystic-server
//	go run ./cmd/mystic-server -auth
//	go run ./cmd/mystic-server -tls -auth  # run gencerts first
package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	pb "mystic-square-grpc/gen/mysticsquare"
	"mystic-square-grpc/internal/auth"
	"mystic-square-grpc/internal/middleware"
	"mystic-square-grpc/internal/server"
	"mystic-square-grpc/internal/store"
	"mystic-square-grpc/internal/tlsconf"
)

func main() {
	addr := flag.String("addr", ":50051", "gRPC listen address")
	useTLS := flag.Bool("tls", false, "require mTLS")
	certDir := flag.String("cert-dir", "certs", "certificate directory")
	useAuth := flag.Bool("auth", false, "require a JWT in authorization metadata")
	secret := flag.String("jwt-secret", envOr("JWT_SECRET", auth.DefaultDemoSecret), "JWT HMAC secret")
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	unary := []grpc.UnaryServerInterceptor{middleware.UnaryRecovery(log), middleware.UnaryLogging(log)}
	stream := []grpc.StreamServerInterceptor{middleware.StreamRecovery(log), middleware.StreamLogging(log)}
	if *useAuth {
		authenticator := auth.NewAuthenticator([]byte(*secret), auth.DefaultPolicy())
		unary = append(unary, authenticator.Unary())
		stream = append(stream, authenticator.Stream())
	}
	opts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(unary...), grpc.ChainStreamInterceptor(stream...)}
	if *useTLS {
		config, err := tlsconf.ServerConfig(*certDir)
		if err != nil {
			log.Error("TLS configuration", "err", err)
			os.Exit(1)
		}
		opts = append(opts, grpc.Creds(credentials.NewTLS(config)))
	}

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Error("listen", "addr", *addr, "err", err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMysticSquareServiceServer(grpcServer, server.New(store.New()))
	reflection.Register(grpcServer) // convenient for grpcurl; disable in production.

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serveErr := make(chan error, 1)
	go func() { serveErr <- grpcServer.Serve(listener) }()
	log.Info("Mystic Square gRPC server started", "addr", listener.Addr(), "mtls", *useTLS, "auth", *useAuth)
	select {
	case err := <-serveErr:
		log.Error("serve", "err", err)
	case <-ctx.Done():
		log.Info("shutting down")
	}
	stopped := make(chan struct{})
	go func() { grpcServer.GracefulStop(); close(stopped) }()
	select {
	case <-stopped:
	case <-time.After(5 * time.Second):
		grpcServer.Stop()
	}
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
