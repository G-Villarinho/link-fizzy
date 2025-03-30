package services

import (
	"context"
	"errors"
	"testing"

	"github.com/g-villarinho/link-fizz-api/mocks"
	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLink(t *testing.T) {
	t.Run("when creation is successful, it should not return an error", func(t *testing.T) {
		mockUtils := new(mocks.UtilsServiceMock)
		mockRepo := new(mocks.LinkRepositoryMock)

		service := &linkService{
			us: mockUtils,
			lr: mockRepo,
		}

		ctx := context.Background()
		expectedShortCode := "abcd1234"

		mockUtils.On("GenerateShortCode", shortCodeLength).Return(expectedShortCode, nil)
		mockRepo.On("CreateLink", mock.Anything, mock.Anything).Return(nil)

		err := service.CreateLink(ctx, "https://example.com")

		assert.NoError(t, err)
		mockUtils.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("when failing to generate the short code, it should return an error", func(t *testing.T) {
		mockUtils := new(mocks.UtilsServiceMock)
		mockRepo := new(mocks.LinkRepositoryMock)

		service := &linkService{
			us: mockUtils,
			lr: mockRepo,
		}

		ctx := context.Background()
		mockUtils.On("GenerateShortCode", shortCodeLength).Return("", errors.New("failed to generate short code"))

		err := service.CreateLink(ctx, "https://example.com")

		assert.Error(t, err)
		assert.Equal(t, "generate short code: failed to generate short code", err.Error())
		mockUtils.AssertExpectations(t)
	})

	t.Run("when failing to create the link in the repository, it should return an error", func(t *testing.T) {
		mockUtils := new(mocks.UtilsServiceMock)
		mockRepo := new(mocks.LinkRepositoryMock)

		service := &linkService{
			us: mockUtils,
			lr: mockRepo,
		}

		ctx := context.Background()
		expectedShortCode := "abcd1234"

		mockUtils.On("GenerateShortCode", shortCodeLength).Return(expectedShortCode, nil)
		mockRepo.On("CreateLink", mock.Anything, mock.Anything).Return(errors.New("repository error"))

		err := service.CreateLink(ctx, "https://example.com")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create link: repository error")
		mockUtils.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetOriginalURLByShortCode(t *testing.T) {
	t.Run("when the short code is found, it should return the original URL", func(t *testing.T) {
		mockRepo := new(mocks.LinkRepositoryMock)

		service := &linkService{
			lr: mockRepo,
		}

		ctx := context.Background()
		shortCode := "abcd1234"
		expectedURL := "https://example.com"

		mockRepo.On("GetOriginalURLByShortCode", ctx, shortCode).Return(expectedURL, nil)

		url, err := service.GetOriginalURLByShortCode(ctx, shortCode)

		assert.NoError(t, err)
		assert.Equal(t, expectedURL, url)
		mockRepo.AssertExpectations(t)
	})

	t.Run("when the short code is not found, it should return ErrLinkNotFound", func(t *testing.T) {
		mockRepo := new(mocks.LinkRepositoryMock)

		service := &linkService{
			lr: mockRepo,
		}

		ctx := context.Background()
		shortCode := "nonexistent"

		mockRepo.On("GetOriginalURLByShortCode", ctx, shortCode).Return("", nil)

		url, err := service.GetOriginalURLByShortCode(ctx, shortCode)

		assert.Error(t, err)
		assert.Equal(t, models.ErrLinkNotFound, err)
		assert.Empty(t, url)
		mockRepo.AssertExpectations(t)
	})

	t.Run("when there is an error in the repository, it should return an error", func(t *testing.T) {
		mockRepo := new(mocks.LinkRepositoryMock)

		service := &linkService{
			lr: mockRepo,
		}

		ctx := context.Background()
		shortCode := "errorcase"

		mockRepo.On("GetOriginalURLByShortCode", ctx, shortCode).Return("", errors.New("database error"))

		url, err := service.GetOriginalURLByShortCode(ctx, shortCode)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get original URL by short code: database error")
		assert.Empty(t, url)
		mockRepo.AssertExpectations(t)
	})
}
