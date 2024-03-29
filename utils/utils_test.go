package utils

import "testing"

const (
	testCert = `
-----BEGIN CERTIFICATE-----
MIIFbTCCA1WgAwIBAgIIBPcug4sZ61UwDQYJKoZIhvcNAQELBQAwTzELMAkGA1UE
BhMCQ04xDTALBgNVBAoTBFggQ0ExGjAYBgNVBAsTEXd3dy54aWV4aWFuYmluLmNu
MRUwEwYDVQQDEwxYIFRMUyBDQSAxQzEwHhcNMjIwOTEyMTE0NzUyWhcNMjQxMjE1
MTE0NzUyWjAYMRYwFAYDVQQDEw14aWV4aWFuYmluLmNuMIICIjANBgkqhkiG9w0B
AQEFAAOCAg8AMIICCgKCAgEAx8+QBAK6I23Dmbpz1BTp7YzjEj2oaFQ7MKt8yqws
iwrZi4d4MNcZvXjEAeDqsSrkJaiOSwgTPrRJQdwAM/q7symttecczdl24wx8jRxB
s6+YuOAzi/wM/F44DETxTxaSJE/zYACyIAagmIAKU8OQObIthBOh5F6uyj++FNBC
z4s1gbNxAC27BiO1oFHmfXjWGwQiGh2gX9kh1Oa2vvsqF9x6OrxdIw7lFkwQfUro
c4dT3nQPtEai3Pz6tbXEWVoKiPjucUTWnm+LH+bwucWxskZNoIYAqOg/ECQ6blON
c+GmZwjgpl0BfDlQnrLXHBaiNfg6rxNKv7IFw1jHpDRoTJCkjTBbEvecE5FANzKI
b+DyLjnWPtl3LkFjTtoRjcH2GWVCPeRVJ2CMFAmPEaBy7JAtcALiMJOBYeiVUPhx
NWcQPZ5YdgfMZklRGJwEBTbdhJ2zJTHqNt0ZU7x46bkLkyUmELydOh/sXkPrpAMu
wm+VEq82ac7DPGKWOjoJRh0cFWzajCN/kywXhKru1i7LslPtXOJ7grlEa6kHXIhf
z8pfMy6QKdv9fLI+KHYbG4XHzCTZMrvAtVxSGDUPQ0t4Np6KlJKqQrhTvCgcgzSG
FXwHkLUHOEJb7OUMQ5L+mZFtKPwSmE4CT7CUiAuV+7i4I96cZiUXBio7FHcnysbP
mgsCAwEAAaOBgzCBgDAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUH
AwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgwFoAUI4WbnJLld/nO
+KkxgpjJeTNLuiUwIAYDVR0RBBkwF4IPKi54aWV4aWFuYmluLmNuhwQBAQEBMA0G
CSqGSIb3DQEBCwUAA4ICAQCA5rDvMBSX0jjayeDmaLWL8RxQlysUo510zS57RP0h
aTRrYArNFV9t3DHWWGHImndqU35TuTQIuqtJbqc4vSYGRkC4Ipxv3/Hu91Uym873
TM3NOeoLbfC7LGopJNYTpaHFbsH58jvSvzW1aGUtbjwyp3fyBZ5Xou/dz0wBbKg0
JldszRIgAgz5T5xxO3fCJnX3DPsH1lONkxgPoHmHnTTz+9pyMI/xMm62A7quqkyk
gKYNX3fCig2t7ujr4Z/wDvncXTbPu2f+6bkCrJlNJRq1rNenViRESNt0oolIytT0
ZzH4meaKvLrPfH64E2iOnNT22uo/Mvves7YLjjVf1K5q6NyaAL/9GjnJtTI2cHZu
dT1WhZNMl2zleBqjlSrXZ2tgAuPsVEWkt5ZJtD67rgZDw/XdkZhDPkUUtB3vseuJ
PyB28hI8Nf3ZCu7KOn1ySunIhcUGa1DFrFUZK2UKFYPOk4W4SNzyaEUheYBQElmv
DHhFBTLzu/3iVyoQg5lE6MVePWQtksvk8q+CpbDnqjLv8FvQw4hIZvgTjCubYh7d
NEEx9HfMlQwug3BSBE1tckb1bjz5mnOZGF1Y8eoZO0dXVrjHifNpgygs0KKzGN+h
cMvoxXmgxg8fmL/Jawn9mrMngucvafJz7nJc5GmJ/G3wEJrSwzXvD72e7TsRnaDf
Pw==
-----END CERTIFICATE-----
`
)

func TestParseCert(t *testing.T) {
	if cert, err := ParseCert([]byte(testCert)); err != nil {
		t.Fatal("parse cert error", err)
	} else {
		t.Logf("DNSNames: %s, IPAddresses: %s, NotAfter: %s, RawSubject: %s",
			cert.DNSNames, cert.IPAddresses, cert.NotAfter, cert.Subject)
	}
}

func TestIsMatch(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"a", "c", "b"}
	if IsMatch(a, b) == false {
		t.Fatal("wrong match")
	}
}
