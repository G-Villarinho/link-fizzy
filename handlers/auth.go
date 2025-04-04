package handlers

import (
	"log/slog"
	"net/http"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/pkgs/requestcontext"
	"github.com/g-villarinho/link-fizz-api/responses"
	"github.com/g-villarinho/link-fizz-api/services"
	jsoniter "github.com/json-iterator/go"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	i  *di.Injector
	as services.AuthService
	rc requestcontext.RequestContext
}

func NewAuthHandler(i *di.Injector) (AuthHandler, error) {
	authService, err := di.Invoke[services.AuthService](i)
	if err != nil {
		return nil, err
	}

	requestContext, err := di.Invoke[requestcontext.RequestContext](i)
	if err != nil {
		return nil, err
	}

	return &authHandler{
		i:  i,
		as: authService,
		rc: requestContext,
	}, nil
}

func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"handler", "login",
		"method", "Login",
	)

	var payload models.LoginPayload
	if err := jsoniter.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Error("decode payload", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := a.as.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		if err == models.ErrUserNotFound {
			logger.Error("invalid credentials", slog.String("error", err.Error()))
			responses.NoContent(w, http.StatusNotFound)
			return
		}

		if err == models.ErrInvalidCredentials {
			logger.Error("invalid credentials", slog.String("error", err.Error()))
			responses.NoContent(w, http.StatusUnauthorized)
			return
		}

		logger.Error("login error", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusInternalServerError)
		return
	}

	responses.JSON(w, http.StatusOK, response)
}

func (a *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"handler", "logout",
		"method", "Logout",
	)

	userID, found := a.rc.GetUserID(r.Context())
	if !found {
		logger.Error("user ID not found in context")
		responses.NoContent(w, http.StatusUnauthorized)
		return
	}

	token, found := a.rc.GetToken(r.Context())
	if !found {
		logger.Error("token not found in context")
		responses.NoContent(w, http.StatusUnauthorized)
		return
	}

	if err := a.as.Logout(r.Context(), userID, token); err != nil {
		if err == models.ErrLogoutAlreadyExists {
			logger.Error("logout already exists", slog.String("error", err.Error()))
			responses.NoContent(w, http.StatusUnauthorized)
			return
		}

		logger.Error("logout error", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusInternalServerError)
		return
	}

	responses.NoContent(w, http.StatusNoContent)
}
