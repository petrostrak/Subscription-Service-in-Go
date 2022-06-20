package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) Routes() http.Handler {
	// create router
	mux := chi.NewRouter()

	// setup middleware
	mux.Use(middleware.Recoverer)

	// define application routes
	mux.Get("/", app.HomePage)

	return mux
}
