package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlinaDovbysheva/go-short-url/internal/app"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/handlers"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/storage"
)

func main() {
	c := app.Config{}
	c.ConfigServerEnv()

	var db storage.DBurl

	if app.DatabaseDsn != "" {
		db = storage.NewInPostgre()
		fmt.Printf("server start on %s  db %s\n", app.ServerURL, app.DatabaseDsn)
		defer db.Close()

	} else if app.FilePath != "" {
		db = storage.NewInFile(app.FilePath)
		fmt.Printf("server start on %s store in file %s\n", app.ServerURL, app.FilePath)
	} else {
		db = storage.NewInMap()
		fmt.Printf("server start on %s store in memory", app.ServerURL)
	}

	h := handlers.NewHandlerServer(db)

	log.Fatal(http.ListenAndServe(app.ServerURL, handlers.GzipHandle(h.Chi)))
}
