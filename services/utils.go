package services

import (
	"math/rand"
	"time"

	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type UtilsService interface {
	GenerateShortCode(length int) (string, error)
}

type utilsService struct {
	i *di.Injector
}

func NewUtilsService(i *di.Injector) (UtilsService, error) {
	return &utilsService{
		i: i,
	}, nil
}

func (u *utilsService) GenerateShortCode(length int) (string, error) {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b), nil
}
