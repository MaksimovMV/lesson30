package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"lesson30/pkg/service"
	"lesson30/pkg/user"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	srv := service.Service{Store: make(map[int]*user.User)}

	r.Use(middleware.Logger)

	r.Post("/create", srv.Create)
	r.Post("/make_friends", srv.MakeFriends)
	r.Delete("/user", srv.DeleteUser)
	r.Get("/friends/{targetID}", srv.GetFriends)
	r.Put("/{targetID}", srv.PutNewAge)

	err := http.ListenAndServe("localhost:8080", r)
	log.Fatal(err)
}
