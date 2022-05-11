package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"lesson30/internal/usecase"
)

func Build(r *chi.Mux, uc *usecase.UseCase) {

	r.Use(middleware.Logger)

	r.Post("/create", uc.Create)
	r.Post("/make_friends", uc.Make)
	r.Delete("/user", uc.User)
	r.Get("/friends/{targetID}", uc.Friends)
	r.Put("/{targetID}", uc.NewAge)
}
