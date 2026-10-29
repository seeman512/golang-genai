// Package middleware contains logging and panic-recovery gRPC interceptors.
package middleware

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func UnaryRecovery(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response any, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error("panic recovered", "method", info.FullMethod, "panic", recovered, "stack", string(debug.Stack()))
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

func StreamRecovery(log *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error("panic recovered", "method", info.FullMethod, "panic", recovered, "stack", string(debug.Stack()))
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(srv, stream)
	}
}

func UnaryLogging(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		started := time.Now()
		response, err := handler(ctx, req)
		logCall(log, ctx, "unary", info.FullMethod, started, err)
		return response, err
	}
}

func StreamLogging(log *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		started := time.Now()
		err := handler(srv, stream)
		logCall(log, stream.Context(), "stream", info.FullMethod, started, err)
		return err
	}
}

func logCall(log *slog.Logger, ctx context.Context, kind, method string, started time.Time, err error) {
	remote := ""
	if p, ok := peer.FromContext(ctx); ok {
		remote = p.Addr.String()
	}
	log.Info("rpc", "kind", kind, "method", method, "code", status.Code(err).String(), "duration", time.Since(started).Round(time.Microsecond), "remote", remote)
}
