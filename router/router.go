package router

import (
	"net/http"

	"github.com/g-villarinho/link-fizz-api/handlers/middlewares"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/router/routes"
	"github.com/gorilla/mux"
)

func SetupRoutes(i *di.Injector) http.Handler {
	r := mux.NewRouter()
	return middlewares.Cors(routes.ConfigureRoutes(r, i))
}
