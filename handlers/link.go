package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/pkgs/requestcontext"
	"github.com/g-villarinho/link-fizz-api/responses"
	"github.com/g-villarinho/link-fizz-api/services"
	jsoniter "github.com/json-iterator/go"
)

type LinkHandler interface {
	CreateLink(w http.ResponseWriter, r *http.Request)
	RedirectLink(w http.ResponseWriter, r *http.Request)
	GetShortURLs(w http.ResponseWriter, r *http.Request)
}

type linkHandler struct {
	i  *di.Injector
	ls services.LinkService
	rs services.RedirectService
	rc requestcontext.RequestContext
}

func NewLinkHandler(i *di.Injector) (LinkHandler, error) {
	linkService, err := di.Invoke[services.LinkService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.LinkService: %w", err)
	}

	redirectService, err := di.Invoke[services.RedirectService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.RedirectService: %w", err)
	}

	requestContext, err := di.Invoke[requestcontext.RequestContext](i)
	if err != nil {
		return nil, fmt.Errorf("invoke requestcontext.RequestContext: %w", err)
	}

	return &linkHandler{
		i:  i,
		ls: linkService,
		rs: redirectService,
		rc: requestContext,
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

	if payload.DestinationURL == "" {
		logger.Error("empty original URL")
		responses.NoContent(w, http.StatusBadRequest)
		return
	}

	userID, found := l.rc.GetUserID(r.Context())
	if !found {
		logger.Error("user ID not found in context")
		responses.NoContent(w, http.StatusUnauthorized)
		return
	}

	if _, err := url.ParseRequestURI(payload.DestinationURL); err != nil {
		responses.NoContent(w, http.StatusBadRequest)
		return
	}

	if err := l.ls.CreateLink(r.Context(), userID, payload.DestinationURL, payload.Title, payload.CustomCode); err != nil {
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

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	userAgent := r.UserAgent()

	originalURL, err := l.rs.GetOriginalURLWithTracking(r.Context(), shortCode, userAgent, ip)
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

func (l *linkHandler) GetShortURLs(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"handler", "link",
		"method", "GetShortURLs",
	)

	userID, found := l.rc.GetUserID(r.Context())
	if !found {
		logger.Error("user ID not found in context")
		responses.NoContent(w, http.StatusUnauthorized)
		return
	}

	shortURLs, err := l.ls.GetUsersShortURLs(r.Context(), userID)
	if err != nil {
		logger.Error("get short URLs", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusInternalServerError)
		return
	}

	responses.JSON(w, http.StatusOK, shortURLs)
}
