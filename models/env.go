package models

import "time"

type Environment struct {
	APIPort        string
	APIURL         string
	DBUser         string
	DBPassword     string
	DBHost         string
	DBPort         string
	DBName         string
	RequestTimeout time.Duration
	Key            Key
}

type Key struct {
	PrivateKey string
	PublicKey  string
}
