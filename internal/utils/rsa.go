package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// ParsePublicKey untuk membaca public key dari file
func ParsePublicKey() (*rsa.PublicKey, error) {
	pubPEM, err := os.ReadFile("internal/keys/public.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %v", err)
	}

	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	return pubKey.(*rsa.PublicKey), nil
}

// ParsePrivateKey untuk membaca private key dari file
func ParsePrivateKey() (*rsa.PrivateKey, error) {
	privPEM, err := os.ReadFile("internal/keys/private.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return privKey, nil
}

// EncryptRSA untuk enkripsi password menggunakan public key
func EncryptRSA(pubKey *rsa.PublicKey, plaintext string) (string, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(plaintext))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt password: %v", err)
	}
	return string(encrypted), nil
}

// DecryptRSA untuk dekripsi password menggunakan private key
func DecryptRSA(privKey *rsa.PrivateKey, ciphertext string) (string, error) {
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, []byte(ciphertext))
	if err != nil {
		return "", fmt.Errorf("failed to decrypt password: %v", err)
	}
	return string(decrypted), nil
}
