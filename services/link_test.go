package services

import (
	"context"
	"errors"
	"testing"

	"github.com/g-villarinho/link-fizz-api/config"
	"github.com/g-villarinho/link-fizz-api/mocks"
	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/google/uuid"
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
		userID := uuid.New().String()

		mockUtils.On("GenerateShortCode", shortCodeLength).Return(expectedShortCode, nil)
		mockRepo.On("CreateLink", mock.Anything, mock.Anything).Return(nil)

		err := service.CreateLink(ctx, userID, "https://example.com")

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
		userID := uuid.New().String()
		err := service.CreateLink(ctx, userID, "https://example.com")

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

		userID := uuid.New().String()
		err := service.CreateLink(ctx, userID, "https://example.com")

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

func TestGetShortURLs(t *testing.T) {
	config.Env.APIURL = "https://api.example.com"

	t.Run("should return short URLs when short codes exist", func(t *testing.T) {
		mockRepo := new(mocks.LinkRepositoryMock)
		service := &linkService{lr: mockRepo}

		ctx := context.Background()
		expectedShortCodes := []string{"abc123", "def456", "ghi789"}
		expectedURLs := []string{
			"https://api.example.com/abc123",
			"https://api.example.com/def456",
			"https://api.example.com/ghi789",
		}

		mockRepo.On("GetAllShortCodes", ctx).Return(expectedShortCodes, nil)

		userID := uuid.New().String()
		result, err := service.GetUsersShortURLs(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedURLs, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return empty slice when no short codes exist", func(t *testing.T) {
		mockRepo := new(mocks.LinkRepositoryMock)
		service := &linkService{lr: mockRepo}

		ctx := context.Background()
		emptyShortCodes := []string{}

		mockRepo.On("GetAllShortCodes", ctx).Return(emptyShortCodes, nil)

		userID := uuid.New().String()
		result, err := service.GetUsersShortURLs(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, []string{}, result)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := new(mocks.LinkRepositoryMock)
		service := &linkService{lr: mockRepo}

		ctx := context.Background()
		expectedError := errors.New("database error")

		mockRepo.On("GetAllShortCodes", ctx).Return(nil, expectedError)

		userID := uuid.New().String()
		result, err := service.GetUsersShortURLs(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get all short codes: database error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should use configured API URL correctly", func(t *testing.T) {
		mockRepo := new(mocks.LinkRepositoryMock)
		service := &linkService{lr: mockRepo}

		originalAPIURL := config.Env.APIURL
		config.Env.APIURL = "https://test.example.com"
		defer func() { config.Env.APIURL = originalAPIURL }()

		ctx := context.Background()
		shortCodes := []string{"test123"}
		expectedURL := "https://test.example.com/test123"

		mockRepo.On("GetAllShortCodes", ctx).Return(shortCodes, nil)

		userID := uuid.New().String()
		result, err := service.GetUsersShortURLs(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, []string{expectedURL}, result)
		mockRepo.AssertExpectations(t)
	})
}
