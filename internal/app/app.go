package app

import (
	"github.com/go-chi/chi/v5"
	"lesson30/internal/controller"
	"lesson30/internal/storage"
	"lesson30/internal/usecase"
	"log"
	"net/http"
)

func Run() {
	r := chi.NewRouter()
	s := storage.NewStorage()
	uc := usecase.NewUseCase(s)
	controller.Build(r, &uc)
	err := http.ListenAndServe("localhost:8080", r)
	log.Fatal(err)
}
