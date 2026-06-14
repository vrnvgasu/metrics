// Package crypto реализует асимметричное шифрование RSA.
package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("crypto.LoadPublicKey: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("crypto.LoadPublicKey: no PEM block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("crypto.LoadPublicKey: %w", err)
	}
	pub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("crypto.LoadPublicKey: not an RSA public key")
	}
	return pub, nil
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("crypto.LoadPrivateKey: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("crypto.LoadPrivateKey: no PEM block")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("crypto.LoadPrivateKey: %w", err)
	}
	return priv, nil
}

func Encrypt(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

func Decrypt(priv *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, priv, data)
}
