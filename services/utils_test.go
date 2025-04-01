package services

import (
	"testing"

	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/stretchr/testify/assert"
)

func TestGenerateShortCode(t *testing.T) {
	t.Run("should generate a short code of the correct length", func(t *testing.T) {
		mockInjector := &di.Injector{}
		service, err := NewUtilsService(mockInjector)

		assert.NoError(t, err)

		length := 8
		shortCode, err := service.GenerateShortCode(length)

		assert.NoError(t, err)
		assert.Len(t, shortCode, length)
	})

	t.Run("should generate a valid short code using the charset", func(t *testing.T) {
		mockInjector := &di.Injector{}
		service, err := NewUtilsService(mockInjector)

		assert.NoError(t, err)

		length := 10
		shortCode, err := service.GenerateShortCode(length)

		assert.NoError(t, err)

		for _, char := range shortCode {
			assert.Contains(t, charset, string(char))
		}
	})
}
