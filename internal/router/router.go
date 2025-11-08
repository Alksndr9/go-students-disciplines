package router

import "github.com/go-chi/chi"

func NewRouter() *chi.Mux {
	router := chi.NewRouter()
	return router
}
