package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

	Env.DBUser = os.Getenv("DB_USER")
	Env.DBPassword = os.Getenv("DB_PASSWORD")
	Env.DBHost = os.Getenv("DB_HOST")
	Env.DBPort = os.Getenv("DB_PORT")
	Env.DBName = os.Getenv("DB_NAME")

	timeoutStr := os.Getenv("REQUEST_TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "30s"
	}
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return fmt.Errorf("invalid timeout format: %w", err)
	}
	Env.RequestTimeout = timeout

	if Env.Key.PrivateKey == "" || Env.Key.PublicKey == "" {
		privateKey, err := LoadKeyFromFile(os.Getenv("KEY_ECDSA_PRIVATE"))
		if err != nil {
			return fmt.Errorf("load private key: %w", err)
		}

		publicKey, err := LoadKeyFromFile(os.Getenv("KEY_ECDSA_PUBLIC"))
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
