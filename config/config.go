package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	PrivateKey []byte
	PublicKey  []byte
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var errkey error
	PrivateKey, errkey = os.ReadFile("internal/keys/private.pem")
	if errkey != nil {
		log.Fatal("Failed to load private key:", errkey)
	}

	PublicKey, errkey = os.ReadFile("internal/keys/public.pem")
	if errkey != nil {
		log.Fatal("Failed to load public key:", errkey)
	}
}
