package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"inventoryseal/internal/config"
	"inventoryseal/internal/httpapi"
	"inventoryseal/internal/service"
	"inventoryseal/internal/store"
)

func main() {
	settings := config.Parse(os.Args[1:])
	db, err := store.Open(settings.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	server := httpapi.New(service.New(db))
	fmt.Printf("inventoryseal listening on %s using %s\n", settings.Listen, settings.DBPath)
	if err := http.ListenAndServe(settings.Listen, server.Handler()); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	_ = context.Background()
}
