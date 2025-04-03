package handlers

import (
	"log/slog"
	"net/http"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/responses"
	"github.com/g-villarinho/link-fizz-api/services"
	jsoniter "github.com/json-iterator/go"
)

type LoginHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type loginHandler struct {
	i  *di.Injector
	ls services.LoginService
}

func NewLoginHandler(i *di.Injector) (LoginHandler, error) {
	ls, err := di.Invoke[services.LoginService](i)
	if err != nil {
		return nil, err
	}

	return &loginHandler{
		i:  i,
		ls: ls,
	}, nil
}

func (l *loginHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	loginResponse, err := l.ls.Login(r.Context(), payload.Email, payload.Password)
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

	responses.JSON(w, http.StatusOK, loginResponse)
}
