package main

import (
	"net/http"

	"github.com/dzivdzi/bookings/pkg/config"
	"github.com/dzivdzi/bookings/pkg/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// Middleware - allows to process a req as it comes into the app and performs an action on it
	mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	// fileServer -> A place to go get static files from
	fileServer := http.FileServer(http.Dir("./static/"))
	// Handle -> StripPrefix Takes the URL go gets and modifies it to something GO knows how to handle, strip the stat and use the fileserver var
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
