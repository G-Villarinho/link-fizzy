package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Env Environment

func LoadEnv() {

	if err := godotenv.Load(".env.development"); err != nil {
		log.Fatal("Error loading .env.development file")
	}

	Env.APIURL = os.Getenv("API_URL")
	Env.APIPort = os.Getenv("API_PORT")
}
