// gateway exposes selected Mystic Square RPCs as REST/JSON using grpc-gateway.
// The gateway forwards Authorization metadata to the upstream gRPC server.
package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	pb "mystic-square-grpc/gen/mysticsquare"
	"mystic-square-grpc/internal/tlsconf"
)

func main() {
	httpAddr := flag.String("http-addr", ":8080", "REST listen address")
	grpcAddr := flag.String("grpc-addr", "localhost:50051", "upstream gRPC address")
	useTLS := flag.Bool("tls", false, "connect to upstream using mTLS")
	certDir := flag.String("cert-dir", "certs", "certificate directory")
	clientCert := flag.String("client-cert", "client", "gateway client certificate")
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var transport credentials.TransportCredentials = insecure.NewCredentials()
	if *useTLS {
		config, err := tlsconf.ClientConfig(*certDir, *clientCert)
		if err != nil {
			log.Error("TLS configuration", "err", err)
			os.Exit(1)
		}
		transport = credentials.NewTLS(config)
	}
	conn, err := grpc.NewClient(*grpcAddr, grpc.WithTransportCredentials(transport))
	if err != nil {
		log.Error("gRPC client", "err", err)
		os.Exit(1)
	}
	defer conn.Close()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions:   protojson.MarshalOptions{UseProtoNames: true, EmitUnpopulated: true},
			UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
		}),
		runtime.WithMetadata(func(_ context.Context, request *http.Request) metadata.MD {
			if value := request.Header.Get("Authorization"); value != "" {
				return metadata.Pairs("authorization", value)
			}
			return nil
		}),
	)
	if err := pb.RegisterMysticSquareServiceHandler(context.Background(), mux, conn); err != nil {
		log.Error("register gateway handlers", "err", err)
		os.Exit(1)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("HTTP request", "method", r.Method, "path", r.URL.RequestURI())
		mux.ServeHTTP(w, r)
	})
	httpServer := &http.Server{Addr: *httpAddr, Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			_ = httpServer.Close()
		}
	}()
	log.Info("gateway started", "http", *httpAddr, "grpc", *grpcAddr, "mtls", *useTLS, "auth_header", strings.ToLower("Authorization"))
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("HTTP server", "err", err)
		os.Exit(1)
	}
}
