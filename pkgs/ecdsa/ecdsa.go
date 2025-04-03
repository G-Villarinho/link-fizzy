package ecdsa

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

type EcdsaKeyPair interface {
	ParseECDSAPrivateKey(pemKey string) (*ecdsa.PrivateKey, error)
	ParseECDSAPublicKey(pemKey string) (*ecdsa.PublicKey, error)
}

type ecdsaKeyPair struct {
	i *di.Injector
}

func NewEcdsaKeyPair(i *di.Injector) (EcdsaKeyPair, error) {
	return &ecdsaKeyPair{
		i: i,
	}, nil
}

func (e *ecdsaKeyPair) ParseECDSAPrivateKey(pemKey string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("failed to parse EC private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (e *ecdsaKeyPair) ParseECDSAPublicKey(pemKey string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to parse EC public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not a valid ECDSA public key")
	}

	return ecdsaPubKey, nil
}
