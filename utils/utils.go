package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"
	"sort"
)

// ParseCert parse tls Cert
func ParseCert(certByte []byte) (*x509.Certificate, error) {
	certBlock, _ := pem.Decode(certByte)
	if certBlock == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err.Error())
	}
	return cert, nil
}

// IsMatch check two array match
func IsMatch(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)

	return reflect.DeepEqual(a, b)
}
