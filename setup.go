package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/g-villarinho/link-fizz-api/databases"
	"github.com/g-villarinho/link-fizz-api/handlers"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/pkgs/ecdsa"
	"github.com/g-villarinho/link-fizz-api/repositories"
	"github.com/g-villarinho/link-fizz-api/services"
)

func initDeps(i *di.Injector) *sql.DB {

	// Databases
	db, err := databases.InitDB()
	if err != nil {
		log.Fatal("failed to initialize database:", err)
	}

	di.Provide(i, func(i *di.Injector) (*sql.DB, error) {
		return db, nil
	})

	// Config
	di.Provide(i, ecdsa.NewEcdsaKeyPair)

	// Handlers
	di.Provide(i, handlers.NewLinkHandler)
	di.Provide(i, handlers.NewUserHandler)
	di.Provide(i, handlers.NewLoginHandler)

	// Services
	di.Provide(i, services.NewLinkService)
	di.Provide(i, services.NewUtilsService)
	di.Provide(i, services.NewUserService)
	di.Provide(i, services.NewLoginService)
	di.Provide(i, services.NewSecurityService)
	di.Provide(i, services.NewTokenService)

	// Repositories
	di.Provide(i, repositories.NewLinkRepository)
	di.Provide(i, repositories.NewUserRepository)

	return db
}

func setupRoutes(i *di.Injector) *http.ServeMux {
	mux := http.NewServeMux()

	setupLinkRoutes(mux, i)
	setupUserRoutes(mux, i)
	setupLoginRoutes(mux, i)

	return mux
}

func setupLinkRoutes(mux *http.ServeMux, i *di.Injector) *http.ServeMux {

	linkHandler, err := di.Invoke[handlers.LinkHandler](i)
	if err != nil {
		log.Fatal("failed to invoke link handler:", err)
	}

	mux.HandleFunc("POST /links", linkHandler.CreateLink)
	mux.HandleFunc("GET /{shortCode}", linkHandler.RedirectLink)
	mux.HandleFunc("GET /links", linkHandler.GetShortURLs)

	return mux
}

func setupUserRoutes(mux *http.ServeMux, i *di.Injector) *http.ServeMux {
	userHandler, err := di.Invoke[handlers.UserHandler](i)
	if err != nil {
		log.Fatal("failed to invoke user handler:", err)
	}

	mux.HandleFunc("POST /users", userHandler.CreateUser)

	return mux
}

func setupLoginRoutes(mux *http.ServeMux, i *di.Injector) *http.ServeMux {
	loginHandler, err := di.Invoke[handlers.LoginHandler](i)
	if err != nil {
		log.Fatal("failed to invoke login handler:", err)
	}

	mux.HandleFunc("POST /login", loginHandler.Login)

	return mux
}
