package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	pb "mystic-square-grpc/gen/mysticsquare"
	"mystic-square-grpc/internal/server"
	"mystic-square-grpc/internal/store"
)

func main() {
	addr := flag.String("addr", ":50051", "gRPC listen address")
	flag.Parse()
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMysticSquareServiceServer(grpcServer, server.New(store.New()))
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()
	log.Printf("Mystic Square gRPC server listening on %s", *addr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	shutdown := make(chan struct{})
	go func() { grpcServer.GracefulStop(); close(shutdown) }()
	select {
	case <-shutdown:
	case <-time.After(5 * time.Second):
		grpcServer.Stop()
	}
}
