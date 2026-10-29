// mint-token prints a local-demo JWT, for example:
//
//	go run ./cmd/mint-token -role moderator
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"mystic-square-grpc/internal/auth"
)

func main() {
	subject := flag.String("sub", "demo-user", "token subject")
	role := flag.String("role", auth.RoleModerator, "role: player | moderator")
	ttl := flag.Duration("ttl", time.Hour, "token lifetime")
	secret := flag.String("secret", envOr("JWT_SECRET", auth.DefaultDemoSecret), "HMAC secret")
	flag.Parse()
	if *role != auth.RolePlayer && *role != auth.RoleModerator {
		log.Fatalf("unknown role %q", *role)
	}
	token, err := auth.Sign([]byte(*secret), auth.Claims{Subject: *subject, Role: *role, ExpiresAt: time.Now().Add(*ttl).Unix()})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
