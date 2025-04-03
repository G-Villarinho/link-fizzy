package models

type Environment struct {
	APIPort string
	APIURL  string
	Key     Key
}

type Key struct {
	PrivateKey string `env:"KEY_ECDSA_PRIVATE"`
	PublicKey  string `env:"KEY_ECDSA_PUBLIC"`
}
