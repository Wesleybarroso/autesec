package transport

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
)

func CertFingerprint(cert *x509.Certificate) string {
	hash := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(hash[:])
}