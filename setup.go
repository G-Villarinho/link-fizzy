package main

import (
	"database/sql"
	"log"

	"github.com/g-villarinho/link-fizz-api/handlers"
	"github.com/g-villarinho/link-fizz-api/handlers/middlewares"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/pkgs/ecdsa"
	"github.com/g-villarinho/link-fizz-api/pkgs/requestcontext"
	"github.com/g-villarinho/link-fizz-api/repositories"
	"github.com/g-villarinho/link-fizz-api/services"
	"github.com/g-villarinho/link-fizz-api/storages"
)

func initDeps(i *di.Injector) *sql.DB {
	db, err := storages.InitDB()
	if err != nil {
		log.Fatal("failed to initialize database:", err)
	}

	di.Provide(i, func(i *di.Injector) (*sql.DB, error) {
		return db, nil
	})

	// Middlewares
	di.Provide(i, middlewares.NewAuthMiddleware)

	// Config
	di.Provide(i, ecdsa.NewEcdsaKeyPair)
	di.Provide(i, requestcontext.NewRequestContext)

	// Handlers
	di.Provide(i, handlers.NewAuthHandler)
	di.Provide(i, handlers.NewLinkHandler)
	di.Provide(i, handlers.NewUserHandler)

	// Services
	di.Provide(i, services.NewAuthService)
	di.Provide(i, services.NewLinkService)
	di.Provide(i, services.NewUtilsService)
	di.Provide(i, services.NewUserService)
	di.Provide(i, services.NewSecurityService)
	di.Provide(i, services.NewTokenService)
	di.Provide(i, services.NewLogoutService)
	di.Provide(i, services.NewRedirectService)
	di.Provide(i, services.NewLinkVisitService)
	di.Provide(i, services.NewSessionService)

	// Repositories
	di.Provide(i, repositories.NewLinkRepository)
	di.Provide(i, repositories.NewUserRepository)
	di.Provide(i, repositories.NewLogoutRepository)
	di.Provide(i, repositories.NewLinkVisitRepository)
	di.Provide(i, repositories.NewSessionRepository)

	return db
}
