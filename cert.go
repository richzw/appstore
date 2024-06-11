package appstore

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// openssl x509 -inform der -in AppleRootCA-G3.cer -out apple_root.pem
const defaultRootPEM = `
-----BEGIN CERTIFICATE-----
MIICQzCCAcmgAwIBAgIILcX8iNLFS5UwCgYIKoZIzj0EAwMwZzEbMBkGA1UEAwwS
QXBwbGUgUm9vdCBDQSAtIEczMSYwJAYDVQQLDB1BcHBsZSBDZXJ0aWZpY2F0aW9u
IEF1dGhvcml0eTETMBEGA1UECgwKQXBwbGUgSW5jLjELMAkGA1UEBhMCVVMwHhcN
MTQwNDMwMTgxOTA2WhcNMzkwNDMwMTgxOTA2WjBnMRswGQYDVQQDDBJBcHBsZSBS
b290IENBIC0gRzMxJjAkBgNVBAsMHUFwcGxlIENlcnRpZmljYXRpb24gQXV0aG9y
aXR5MRMwEQYDVQQKDApBcHBsZSBJbmMuMQswCQYDVQQGEwJVUzB2MBAGByqGSM49
AgEGBSuBBAAiA2IABJjpLz1AcqTtkyJygRMc3RCV8cWjTnHcFBbZDuWmBSp3ZHtf
TjjTuxxEtX/1H7YyYl3J6YRbTzBPEVoA/VhYDKX1DyxNB0cTddqXl5dvMVztK517
IDvYuVTZXpmkOlEKMaNCMEAwHQYDVR0OBBYEFLuw3qFYM4iapIqZ3r6966/ayySr
MA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgEGMAoGCCqGSM49BAMDA2gA
MGUCMQCD6cHEFl4aXTQY2e3v9GwOAEZLuN+yRhHFD/3meoyhpmvOwgPUnPWTxnS4
at+qIxUCMG1mihDK1A3UT82NQz60imOlM27jbdoXt2QfyFMm+YhidDkLF1vLUagM
6BgD56KyKA==
-----END CERTIFICATE-----
`

type Cert struct {
	rootCertPool *x509.CertPool
}

func newCert(rootCertPool *x509.CertPool) *Cert {
	if rootCertPool == nil {
		rootCertPool = x509.NewCertPool()
		rootCertPool.AppendCertsFromPEM([]byte(defaultRootPEM))
	}
	return &Cert{rootCertPool: rootCertPool}
}

func (c *Cert) parseCert(certStr string) (*x509.Certificate, error) {
	certByte, err := base64.StdEncoding.DecodeString(certStr)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certByte)
}

func (c *Cert) extractPublicKeyFromToken(token string) (*ecdsa.PublicKey, error) {
	headerStr, _, _ := strings.Cut(token, ".")
	headerByte, err := base64.RawStdEncoding.DecodeString(headerStr)
	if err != nil {
		return nil, err
	}

	var header struct {
		Alg string   `json:"alg"`
		X5c []string `json:"x5c"`
	}
	err = json.Unmarshal(headerByte, &header)
	if err != nil {
		return nil, err
	}
	if len(header.X5c) == 0 {
		return nil, errors.New("appstore found no certificates in x5c header field")
	}

	opts := x509.VerifyOptions{Roots: c.rootCertPool}

	leafCert, err := c.parseCert(header.X5c[0])
	if err != nil {
		return nil, fmt.Errorf("appstore failed to parse leaf certificate: %w", err)
	}
	header.X5c = header.X5c[1:]

	pk, ok := leafCert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("appstore public key must be of type ecdsa.PublicKey")
	}

	// Build intermediate cert pool if there is more than 1 certificate in the header
	if len(header.X5c) > 0 {
		opts.Intermediates = x509.NewCertPool()

		for i, certStr := range header.X5c {
			cert, err := c.parseCert(certStr)
			if err != nil {
				return nil, fmt.Errorf("appstore failed to parse intermediate certificate %d: %w", i, err)
			}
			opts.Intermediates.AddCert(cert)
		}
	}

	_, err = leafCert.Verify(opts)
	if err != nil {
		return nil, fmt.Errorf("appstore failed to verify leaf certificate: %w", err)
	}

	return pk, nil
}
