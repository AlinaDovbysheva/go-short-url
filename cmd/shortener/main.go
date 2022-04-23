package main

import (
	"github.com/AlinaDovbysheva/go-short-url/internal/app"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/handlers"
	"log"
	"net/http"
)

func main() {
	appHandlerServ := handlers.NewHandlerServer()
	http.HandleFunc("/", appHandlerServ.HandlerServerMain)
	server := &http.Server{
		Addr: app.ServerUrl,
	}
	log.Fatal(server.ListenAndServe())
}
