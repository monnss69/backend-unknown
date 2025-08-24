package main

import (
	"log"
	"net/http"
	"os"

	"componentstore/internal/components"
	"componentstore/internal/database"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := database.Open(dsn)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	compHandler := components.NewHandler(db)
	mux.HandleFunc("/components", compHandler.Components)
	mux.HandleFunc("/components/", compHandler.ComponentByID)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if port := os.Getenv("PORT"); port != "" {
		srv.Addr = ":" + port
	}

	log.Printf("listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}
