package routes

import (
	"log"
	"net/http"

	"github.com/g-villarinho/link-fizz-api/handlers"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

func GetUserRoutes(i *di.Injector) []Route {
	userHandler, err := di.Invoke[handlers.UserHandler](i)
	if err != nil {
		log.Fatal("failed to inject user handler:", err)
	}

	return []Route{
		{
			Method:         http.MethodGet,
			Path:           "/me",
			Handler:        userHandler.GetProfile,
			AllowAnonymous: false,
		},
		{
			Method:         http.MethodPut,
			Path:           "/users",
			Handler:        userHandler.UpdateUser,
			AllowAnonymous: false,
		},
		{
			Method:         http.MethodDelete,
			Path:           "/users",
			Handler:        userHandler.DeleteUser,
			AllowAnonymous: false,
		},
		{
			Method:         http.MethodPatch,
			Path:           "/users/password",
			Handler:        userHandler.UpdatePassword,
			AllowAnonymous: false,
		},
	}
}
