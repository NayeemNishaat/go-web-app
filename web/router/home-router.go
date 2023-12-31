package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/web/controller"
)

func HomeRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", controller.App.HomePage)

	return mux
}
