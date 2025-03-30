package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/g-villarinho/link-fizz-api/databases"
	"github.com/g-villarinho/link-fizz-api/handlers"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
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

	// Handlers
	di.Provide(i, handlers.NewLinkHandler)

	// Services
	di.Provide(i, services.NewLinkService)
	di.Provide(i, services.NewUtilsService)

	// Repositories
	di.Provide(i, repositories.NewLinkRepository)

	return db
}

func setupRoutes(i *di.Injector) *http.ServeMux {
	mux := http.NewServeMux()

	linkHandler, err := di.Invoke[handlers.LinkHandler](i)
	if err != nil {
		log.Fatal("failed to invoke link handler:", err)
	}

	mux.HandleFunc("POST /links", linkHandler.CreateLink)
	mux.HandleFunc("GET /{shortCode}", linkHandler.RedirectLink)

	return mux
}
