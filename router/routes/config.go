package routes

import (
	"log"
	"net/http"

	"github.com/g-villarinho/link-fizz-api/handlers/middlewares"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/gorilla/mux"
)

type Route struct {
	Method         string
	Path           string
	Handler        http.HandlerFunc
	AllowAnonymous bool
}

func ConfigureRoutes(r *mux.Router, i *di.Injector) *mux.Router {
	authMiddleware, err := di.Invoke[middlewares.AuthMiddleware](i)
	if err != nil {
		log.Fatal("Failed to create auth middleware:", err)
	}

	allRoutes := []Route{}
	allRoutes = append(allRoutes, GetAuthRoutes(i)...)
	allRoutes = append(allRoutes, GetUserRoutes(i)...)
	allRoutes = append(allRoutes, GetLinkRoutes(i)...)

	for _, route := range allRoutes {
		var handler http.Handler = route.Handler

		if !route.AllowAnonymous {
			handler = authMiddleware.Authenticate(handler)
		}

		r.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
		}).Methods(route.Method)
	}

	return r
}
