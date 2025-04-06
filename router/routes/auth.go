package routes

import (
	"log"
	"net/http"

	"github.com/g-villarinho/link-fizz-api/handlers"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

func GetAuthRoutes(i *di.Injector) []Route {
	authHandler, err := di.Invoke[handlers.AuthHandler](i)
	if err != nil {
		log.Fatal("failed to create auth handler:", err)
	}

	return []Route{
		{
			Method:         http.MethodPost,
			Path:           "/register",
			Handler:        authHandler.Register,
			AllowAnonymous: true,
		},
		{
			Method:         http.MethodPost,
			Path:           "/login",
			Handler:        authHandler.Login,
			AllowAnonymous: true,
		},
		{
			Method:         http.MethodPost,
			Path:           "/logout",
			Handler:        authHandler.Logout,
			AllowAnonymous: false,
		},
	}
}
