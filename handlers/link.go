package handlers

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/responses"
	"github.com/g-villarinho/link-fizz-api/services"
	jsoniter "github.com/json-iterator/go"
)

type LinkHandler interface {
	CreateLink(w http.ResponseWriter, r *http.Request)
	RedirectLink(w http.ResponseWriter, r *http.Request)
}

type linkHandler struct {
	i  *di.Injector
	ls services.LinkService
}

func NewLinkHandler(i *di.Injector) (LinkHandler, error) {
	ls, err := di.Invoke[services.LinkService](i)
	if err != nil {
		return nil, err
	}

	return &linkHandler{
		i:  i,
		ls: ls,
	}, nil
}

func (l *linkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"handler", "link",
		"method", "CreateLink",
	)

	var payload models.LinkPayload
	if err := jsoniter.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Error("decode payload", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if payload.OriginalURL == "" {
		logger.Error("empty original URL")
		responses.NoContent(w, http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(payload.OriginalURL); err != nil {
		responses.NoContent(w, http.StatusBadRequest)
		return
	}

	if err := l.ls.CreateLink(r.Context(), payload.OriginalURL); err != nil {
		logger.Error("create link", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusInternalServerError)
		return
	}

	responses.NoContent(w, http.StatusCreated)
}

func (l *linkHandler) RedirectLink(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"handler", "link",
		"method", "RedirectLink",
	)

	shortCode := r.URL.Path[1:]
	if shortCode == "" {
		logger.Error("empty short code")
		responses.NoContent(w, http.StatusBadRequest)
		return
	}

	originalURL, err := l.ls.GetOriginalURLByShortCode(r.Context(), shortCode)
	if err != nil {
		if err == models.ErrLinkNotFound {
			logger.Error("original URL not found")
			responses.NoContent(w, http.StatusNotFound)
			return
		}

		logger.Error("get original URL by short code", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
