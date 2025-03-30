package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/g-villarinho/link-fizz-api/mocks"
	"github.com/g-villarinho/link-fizz-api/models"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLink(t *testing.T) {
	t.Run("should return 201 when link is created successfully", func(t *testing.T) {
		mockLinkService := new(mocks.LinkServiceMock)
		handler := &linkHandler{
			ls: mockLinkService,
		}

		originalURL := "https://example.com"
		payload := models.LinkPayload{OriginalURL: originalURL}

		mockLinkService.On("CreateLink", mock.Anything, originalURL).Return(nil)

		payloadBytes, _ := jsoniter.Marshal(payload)
		req, err := http.NewRequest("POST", "/link", bytes.NewReader(payloadBytes))
		assert.NoError(t, err)
		resp := httptest.NewRecorder()

		handler.CreateLink(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockLinkService.AssertExpectations(t)
	})

	t.Run("should return 400 when the payload is invalid", func(t *testing.T) {
		mockLinkService := new(mocks.LinkServiceMock)
		handler := &linkHandler{
			ls: mockLinkService,
		}

		invalidPayload := "{invalid_json}"
		req, err := http.NewRequest("POST", "/link", bytes.NewReader([]byte(invalidPayload)))
		assert.NoError(t, err)
		resp := httptest.NewRecorder()

		handler.CreateLink(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return 400 when the original URL is empty", func(t *testing.T) {
		mockLinkService := new(mocks.LinkServiceMock)
		handler := &linkHandler{
			ls: mockLinkService,
		}

		invalidPayload := models.LinkPayload{OriginalURL: ""}
		payloadBytes, _ := jsoniter.Marshal(invalidPayload)
		req, err := http.NewRequest("POST", "/link", bytes.NewReader(payloadBytes))
		assert.NoError(t, err)
		resp := httptest.NewRecorder()

		handler.CreateLink(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return 400 when the original URL is invalid", func(t *testing.T) {
		mockLinkService := new(mocks.LinkServiceMock)
		handler := &linkHandler{
			ls: mockLinkService,
		}

		invalidURLPayload := models.LinkPayload{OriginalURL: "invalid_url"}
		payloadBytes, _ := jsoniter.Marshal(invalidURLPayload)
		req, err := http.NewRequest("POST", "/link", bytes.NewReader(payloadBytes))
		assert.NoError(t, err)
		resp := httptest.NewRecorder()

		handler.CreateLink(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return 500 when an error occurs while creating the link", func(t *testing.T) {
		mockLinkService := new(mocks.LinkServiceMock)
		handler := &linkHandler{
			ls: mockLinkService,
		}

		originalURL := "https://example.com"
		payload := models.LinkPayload{OriginalURL: originalURL}

		mockLinkService.On("CreateLink", mock.Anything, originalURL).Return(errors.New("repository error"))

		payloadBytes, _ := jsoniter.Marshal(payload)
		req, err := http.NewRequest("POST", "/link", bytes.NewReader(payloadBytes))
		assert.NoError(t, err)
		resp := httptest.NewRecorder()

		handler.CreateLink(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockLinkService.AssertExpectations(t)
	})
}

func TestRedirectLink(t *testing.T) {
	t.Run("should return 302 and redirect to the original URL when short code exists", func(t *testing.T) {
		mockLinkService := new(mocks.LinkServiceMock)
		handler := &linkHandler{
			ls: mockLinkService,
		}

		shortCode := "abcd1234"
		expectedURL := "https://example.com"

		mockLinkService.On("GetOriginalURLByShortCode", mock.Anything, shortCode).Return(expectedURL, nil)

		req, err := http.NewRequest("GET", "/"+shortCode, nil)
		assert.NoError(t, err)
		resp := httptest.NewRecorder()

		handler.RedirectLink(resp, req)

		assert.Equal(t, http.StatusFound, resp.Code)
		assert.Equal(t, expectedURL, resp.Header().Get("Location"))
		mockLinkService.AssertExpectations(t)
	})

	t.Run("should return 400 when the short code is empty", func(t *testing.T) {
		mockLinkService := new(mocks.LinkServiceMock)
		handler := &linkHandler{
			ls: mockLinkService,
		}

		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		resp := httptest.NewRecorder()

		handler.RedirectLink(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return 404 when the short code does not exist", func(t *testing.T) {
		mockLinkService := new(mocks.LinkServiceMock)
		handler := &linkHandler{
			ls: mockLinkService,
		}

		shortCode := "nonexistent"
		mockLinkService.On("GetOriginalURLByShortCode", mock.Anything, shortCode).Return("", models.ErrLinkNotFound)

		req, err := http.NewRequest("GET", "/"+shortCode, nil)
		assert.NoError(t, err)
		resp := httptest.NewRecorder()

		handler.RedirectLink(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockLinkService.AssertExpectations(t)
	})

	t.Run("should return 500 when there is an error fetching the original URL", func(t *testing.T) {
		mockLinkService := new(mocks.LinkServiceMock)
		handler := &linkHandler{
			ls: mockLinkService,
		}

		shortCode := "errorcase"
		mockLinkService.On("GetOriginalURLByShortCode", mock.Anything, shortCode).Return("", errors.New("database error"))

		req, err := http.NewRequest("GET", "/"+shortCode, nil)
		assert.NoError(t, err)
		resp := httptest.NewRecorder()

		handler.RedirectLink(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockLinkService.AssertExpectations(t)
	})
}
