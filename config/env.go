package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/joho/godotenv"
)

var Env models.Environment

func LoadEnv() error {

	if err := godotenv.Load(".env.development"); err != nil {
		log.Fatal("Error loading .env.development file")
	}

	Env.APIURL = os.Getenv("API_URL")
	Env.APIPort = os.Getenv("API_PORT")

	if Env.Key.PrivateKey == "" || Env.Key.PublicKey == "" {
		privateKey, err := LoadKeyFromFile("ecdsa_private.pem")
		if err != nil {
			return fmt.Errorf("load private key: %w", err)
		}

		publicKey, err := LoadKeyFromFile("ecdsa_public.pem")
		if err != nil {
			return fmt.Errorf("load public key: %w", err)
		}

		Env.Key.PrivateKey = privateKey
		Env.Key.PublicKey = publicKey
	}

	return nil
}

func LoadKeyFromFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return strings.TrimSpace(string(data)), nil
}
