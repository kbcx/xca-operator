package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

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
