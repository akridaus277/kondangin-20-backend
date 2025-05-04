package main

import (
	"fmt"
	"kondangin-backend/internal/utils"
)

func main() {
	// Parse public key dan private key
	pubKey, err := utils.ParsePublicKey()
	if err != nil {
		fmt.Println("Error parsing public key:", err)
		return
	}
	// Encrypt password
	originalPassword := "supersecret"
	encryptedPassword, err := utils.EncryptRSA(pubKey, originalPassword)
	if err != nil {
		fmt.Println("Error encrypting password:", err)
		return
	}
	fmt.Println("Encrypted Password:", encryptedPassword)
}
