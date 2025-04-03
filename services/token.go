package services

import (
	"context"
	"fmt"
	"time"

	"github.com/g-villarinho/link-fizz-api/config"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/pkgs/ecdsa"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateToken(ctx context.Context, userID string) (string, error)
}

type tokenService struct {
	i  *di.Injector
	kp ecdsa.EcdsaKeyPair
}

func NewTokenService(i *di.Injector) (TokenService, error) {
	ecdsaKeyPair, err := di.Invoke[ecdsa.EcdsaKeyPair](i)
	if err != nil {
		return nil, fmt.Errorf("invoke ecdsa.EcdsaKeyPair: %w", err)
	}

	return &tokenService{
		i:  i,
		kp: ecdsaKeyPair,
	}, nil
}

func (t *tokenService) GenerateToken(ctx context.Context, userID string) (string, error) {
	privateKey, err := t.kp.ParseECDSAPrivateKey(config.Env.Key.PrivateKey)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"iss": "link-fizz-app",
		"sub": userID,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}
