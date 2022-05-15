package main

import (
	"github.com/AlinaDovbysheva/go-short-url/internal/app"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/handlers"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	c := app.Config{}
	c.ConfigServerEnv()

	//db := storage.NewInMap()
	db := storage.NewInFile(app.FilePath)

	h := handlers.NewHandlerServer(db)

	r := chi.NewRouter()
	r.Get("/{id}", h.HandlerServerGet)
	r.Post("/api/shorten", h.HandlerServerPostJSON)
	r.Post("/", h.HandlerServerPost)

	log.Fatal(http.ListenAndServe(app.ServerURL, r))
}
