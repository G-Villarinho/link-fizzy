package routes

import (
	"log"
	"net/http"

	"github.com/g-villarinho/link-fizz-api/handlers"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

func GetLinkRoutes(i *di.Injector) []Route {
	linkHandler, err := di.Invoke[handlers.LinkHandler](i)
	if err != nil {
		log.Fatal("failed to inject link handler:", err)
	}

	return []Route{
		{
			Method:         http.MethodPost,
			Path:           "/links",
			Handler:        linkHandler.CreateLink,
			AllowAnonymous: false,
		},
		{
			Method:         http.MethodGet,
			Path:           "/{shortCode}",
			Handler:        linkHandler.RedirectLink,
			AllowAnonymous: true,
		},
		{
			Method:         http.MethodGet,
			Path:           "/me/links",
			Handler:        linkHandler.GetLinks,
			AllowAnonymous: false,
		},
		{
			Method:         http.MethodGet,
			Path:           "/me/links/{shortCode}",
			Handler:        linkHandler.GetLinkDetails,
			AllowAnonymous: false,
		},
	}
}
