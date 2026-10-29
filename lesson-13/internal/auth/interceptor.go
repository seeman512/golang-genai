package auth

import (
	"context"
	"slices"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "mystic-square-grpc/gen/mysticsquare"
)

type contextKey struct{}

// ClaimsFromContext returns the claims installed by the authentication interceptor.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(contextKey{}).(*Claims)
	return claims, ok
}

// Policy maps fully-qualified gRPC methods to allowed roles.
type Policy map[string][]string

func DefaultPolicy() Policy {
	readers := []string{RolePlayer, RoleModerator}
	moderators := []string{RoleModerator}
	return Policy{
		pb.MysticSquareService_Register_FullMethodName:                moderators,
		pb.MysticSquareService_ConfirmRegistration_FullMethodName:     moderators,
		pb.MysticSquareService_CompleteRegistration_FullMethodName:    moderators,
		pb.MysticSquareService_WatchRegistrationStatus_FullMethodName: readers,
		pb.MysticSquareService_CreateGame_FullMethodName:              readers,
		pb.MysticSquareService_GetGame_FullMethodName:                 readers,
		pb.MysticSquareService_SubmitGameMoves_FullMethodName:         readers,
	}
}

type Authenticator struct {
	secret []byte
	policy Policy
	now    func() time.Time
}

// NewAuthenticator creates a JWT authenticator with the supplied RBAC policy.
func NewAuthenticator(secret []byte, policy Policy) *Authenticator {
	return &Authenticator{secret: secret, policy: policy, now: time.Now}
}

// Unary returns an interceptor for unary RPCs.
func (a *Authenticator) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, err := a.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// Stream returns an interceptor for server-, client-, and bidirectional streams.
func (a *Authenticator) Stream() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, err := a.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		return handler(srv, &wrappedStream{ServerStream: stream, ctx: ctx})
	}
}

func (a *Authenticator) authorize(ctx context.Context, method string) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	raw := values[0]
	if len(raw) >= len("Bearer ") && strings.EqualFold(raw[:len("Bearer ")], "Bearer ") {
		raw = raw[len("Bearer "):]
	}
	claims, err := Verify(a.secret, raw, a.now())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	// Reflection is useful for grpcurl. It still needs a valid token, but does
	// not need a business-method policy entry.
	if !strings.HasPrefix(method, "/grpc.reflection.") {
		allowed, known := a.policy[method]
		if !known || !slices.Contains(allowed, claims.Role) {
			return nil, status.Errorf(codes.PermissionDenied, "role %q is not allowed to call %s", claims.Role, method)
		}
	}
	return context.WithValue(ctx, contextKey{}, claims), nil
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *wrappedStream) Context() context.Context { return s.ctx }

type perRPCToken struct {
	token      string
	requireTLS bool
}

// NewPerRPCToken returns credentials that add a Bearer token to every RPC.
func NewPerRPCToken(token string, requireTLS bool) credentials.PerRPCCredentials {
	return perRPCToken{token: token, requireTLS: requireTLS}
}

func (t perRPCToken) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Bearer " + t.token}, nil
}

func (t perRPCToken) RequireTransportSecurity() bool { return t.requireTLS }
