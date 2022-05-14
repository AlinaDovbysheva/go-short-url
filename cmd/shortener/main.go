package main

import (
	"github.com/AlinaDovbysheva/go-short-url/internal/app"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/handlers"
	"log"
	"net/http"
)

func main() {
	c := app.Config{}
	c.ConfigServerEnv()

	appHandler := handlers.NewHandlerServer()

	server := &http.Server{
		Addr:    app.ServerURL,
		Handler: appHandler,
	}

	log.Fatal(server.ListenAndServe())
}
