// Package tlsconf builds the TLS configurations used by the local mTLS demo.
package tlsconf

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
)

// ServerConfig requires a client certificate signed by the demo CA.
func ServerConfig(dir string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(filepath.Join(dir, "server.pem"), filepath.Join(dir, "server-key.pem"))
	if err != nil {
		return nil, fmt.Errorf("load server certificate (run go run ./cmd/gencerts): %w", err)
	}
	pool, err := loadPool(filepath.Join(dir, "ca.pem"))
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// ClientConfig verifies the server and optionally presents a client certificate.
// certName is client, rogue-client, or none.
func ClientConfig(dir, certName string) (*tls.Config, error) {
	pool, err := loadPool(filepath.Join(dir, "ca.pem"))
	if err != nil {
		return nil, err
	}
	config := &tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12}
	if certName == "none" {
		return config, nil
	}
	cert, err := tls.LoadX509KeyPair(filepath.Join(dir, certName+".pem"), filepath.Join(dir, certName+"-key.pem"))
	if err != nil {
		return nil, fmt.Errorf("load client certificate %q: %w", certName, err)
	}
	config.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) { return &cert, nil }
	return config, nil
}

func loadPool(path string) (*x509.CertPool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read CA %s: %w", path, err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("no certificates found in %s", path)
	}
	return pool, nil
}
