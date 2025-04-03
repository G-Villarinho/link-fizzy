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

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	i  *di.Injector
	ur services.UserService
}

func NewUserHandler(i *di.Injector) (UserHandler, error) {
	ur, err := di.Invoke[services.UserService](i)
	if err != nil {
		return nil, err
	}

	return &userHandler{
		i:  i,
		ur: ur,
	}, nil
}

func (u *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"handler", "user",
		"method", "CreateUser",
	)

	var payload models.CreateUserPayload
	if err := jsoniter.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Error("decode payload", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := u.ur.CreateUser(r.Context(), payload.Name, payload.Email, payload.Password); err != nil {
		if err == models.ErrUserAlreadyExists {
			logger.Error("user already exists", slog.String("error", err.Error()))
			responses.NoContent(w, http.StatusConflict)
			return
		}

		logger.Error("create user", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusInternalServerError)
		return
	}

	responses.NoContent(w, http.StatusCreated)
}
