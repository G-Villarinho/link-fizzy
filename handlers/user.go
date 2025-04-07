package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/pkgs/requestcontext"
	"github.com/g-villarinho/link-fizz-api/responses"
	"github.com/g-villarinho/link-fizz-api/services"
	jsoniter "github.com/json-iterator/go"
)

type UserHandler interface {
	GetProfile(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	UpdatePassword(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	i  *di.Injector
	rc requestcontext.RequestContext
	ur services.UserService
}

func NewUserHandler(i *di.Injector) (UserHandler, error) {
	userService, err := di.Invoke[services.UserService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.UserService: %w", err)
	}

	requestContext, err := di.Invoke[requestcontext.RequestContext](i)
	if err != nil {
		return nil, fmt.Errorf("invoke requestcontext.RequestContext: %w", err)
	}

	return &userHandler{
		i:  i,
		rc: requestContext,
		ur: userService,
	}, nil
}

func (u *userHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"handler", "user",
		"method", "GetProfile",
	)

	userID, found := u.rc.GetUserID(r.Context())
	if !found {
		logger.Error("user ID not found in context")
		responses.NoContent(w, http.StatusUnauthorized)
		return
	}

	resp, err := u.ur.GetUserByID(r.Context(), userID)
	if err != nil {
		logger.Error("get user by ID", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusInternalServerError)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
}

func (u *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"handler", "user",
		"method", "UpdateUser",
	)

	var payload models.UpdateUserPayload
	if err := jsoniter.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Error("decode payload", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID, found := u.rc.GetUserID(r.Context())
	if !found {
		logger.Error("user ID not found in context")
		responses.NoContent(w, http.StatusUnauthorized)
		return
	}

	if err := u.ur.UpdateUser(r.Context(), userID, payload.Name, payload.Email); err != nil {
		if err == models.ErrUserNotFound {
			logger.Error("user not found", slog.String("error", err.Error()))
			responses.NoContent(w, http.StatusNotFound)
			return
		}

		if err == models.ErrUserAlreadyExists {
			logger.Error("user already exists", slog.String("error", err.Error()))
			responses.NoContent(w, http.StatusConflict)
			return
		}

		logger.Error("update user", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusInternalServerError)
		return
	}

	responses.NoContent(w, http.StatusNoContent)
}

func (u *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func (u *userHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"handler", "user",
		"method", "UpdatePassword",
	)

	var payload models.UpdatePasswordPayload
	if err := jsoniter.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Error("decode payload", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID, found := u.rc.GetUserID(r.Context())
	if !found {
		logger.Error("user ID not found in context")
		responses.NoContent(w, http.StatusUnauthorized)
		return
	}

	if err := u.ur.UpdatePassword(r.Context(), userID, payload.CurrentPassword, payload.NewPassword); err != nil {
		if err == models.ErrUserNotFound {
			logger.Error("user not found", slog.String("error", err.Error()))
			responses.NoContent(w, http.StatusUnauthorized)
			return
		}

		if err == models.ErrInvalidCredentials {
			logger.Error("invalid credentials", slog.String("error", err.Error()))
			responses.NoContent(w, http.StatusUnauthorized)
			return
		}

		logger.Error("update password", slog.String("error", err.Error()))
		responses.NoContent(w, http.StatusInternalServerError)
		return
	}

	responses.NoContent(w, http.StatusNoContent)
}
