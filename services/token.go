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
	ValidateToken(ctx context.Context, tokenString string) (string, error)
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

func (t *tokenService) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	publicKey, err := t.kp.ParseECDSAPublicKey(config.Env.Key.PublicKey)
	if err != nil {
		return "", fmt.Errorf("parse public key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", fmt.Errorf("invalid token claims")
		}
		return userID, nil
	}

	return "", fmt.Errorf("invalid token")
}
