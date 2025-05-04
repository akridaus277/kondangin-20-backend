package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// GenerateKeys menghasilkan RSA private dan public keys dan menyimpannya dalam file
func GenerateKeys() error {
	// Generate private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Create private key file
	privFile, err := os.Create("keys/private.pem")
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privFile.Close()

	// Encode private key to PEM format
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	privPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}
	if err := pem.Encode(privFile, privPem); err != nil {
		return fmt.Errorf("failed to encode private key: %v", err)
	}

	// Generate public key
	pub := &priv.PublicKey

	// Create public key file
	pubFile, err := os.Create("keys/public.pem")
	if err != nil {
		return fmt.Errorf("failed to create public key file: %v", err)
	}
	defer pubFile.Close()

	// Marshal and encode the public key to PEM format
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}

	pubPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	if err := pem.Encode(pubFile, pubPem); err != nil {
		return fmt.Errorf("failed to encode public key: %v", err)
	}

	return nil
}
