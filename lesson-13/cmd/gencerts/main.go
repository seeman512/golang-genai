// gencerts creates local-only certificates for the mTLS demonstration.
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

type certificatePair struct {
	cert *x509.Certificate
	key  *ecdsa.PrivateKey
}

func main() {
	dir := flag.String("dir", "certs", "certificate directory")
	force := flag.Bool("force", false, "overwrite existing certificates")
	flag.Parse()
	if _, err := os.Stat(filepath.Join(*dir, "ca.pem")); err == nil && !*force {
		fmt.Printf("certificates already exist in %s (use -force to replace)\n", *dir)
		return
	}
	if err := os.MkdirAll(*dir, 0o755); err != nil {
		log.Fatal(err)
	}
	ca := newCA("Mystic Square Demo CA")
	writePair(*dir, "ca", ca)
	server := newLeaf(ca, "mystic-server", x509.ExtKeyUsageServerAuth, []string{"localhost"}, []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")})
	writePair(*dir, "server", server)
	client := newLeaf(ca, "mystic-client", x509.ExtKeyUsageClientAuth, nil, nil)
	writePair(*dir, "client", client)
	rogueCA := newCA("Untrusted Demo CA")
	writeCert(*dir, "rogue-ca", rogueCA)
	writePair(*dir, "rogue-client", newLeaf(rogueCA, "rogue-client", x509.ExtKeyUsageClientAuth, nil, nil))
	fmt.Printf("certificates written to %s/\n", *dir)
}

func newCA(commonName string) certificatePair {
	key := newKey()
	template := &x509.Certificate{SerialNumber: serial(), Subject: pkix.Name{CommonName: commonName, Organization: []string{"Mystic Square Demo"}}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().AddDate(1, 0, 0), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageCRLSign}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	check(err)
	cert, err := x509.ParseCertificate(der)
	check(err)
	return certificatePair{cert: cert, key: key}
}

func newLeaf(ca certificatePair, commonName string, usage x509.ExtKeyUsage, dnsNames []string, ips []net.IP) certificatePair {
	key := newKey()
	template := &x509.Certificate{SerialNumber: serial(), Subject: pkix.Name{CommonName: commonName, Organization: []string{"Mystic Square Demo"}}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().AddDate(1, 0, 0), KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{usage}, DNSNames: dnsNames, IPAddresses: ips}
	der, err := x509.CreateCertificate(rand.Reader, template, ca.cert, &key.PublicKey, ca.key)
	check(err)
	cert, err := x509.ParseCertificate(der)
	check(err)
	return certificatePair{cert: cert, key: key}
}

func writePair(dir, name string, pair certificatePair) {
	writeCert(dir, name, pair)
	der, err := x509.MarshalECPrivateKey(pair.key)
	check(err)
	writePEM(filepath.Join(dir, name+"-key.pem"), "EC PRIVATE KEY", der, 0o600)
}

func writeCert(dir, name string, pair certificatePair) {
	writePEM(filepath.Join(dir, name+".pem"), "CERTIFICATE", pair.cert.Raw, 0o644)
}
func writePEM(path, kind string, data []byte, mode os.FileMode) {
	check(os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: kind, Bytes: data}), mode))
}
func newKey() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	check(err)
	return key
}
func serial() *big.Int {
	n, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 120))
	check(err)
	return n
}
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
