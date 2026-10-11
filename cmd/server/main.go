package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Bazu-Bazu/Http-Server-On-Go/internal/httpapi"
	"github.com/Bazu-Bazu/Http-Server-On-Go/internal/user"
)

func main() {
	repo := user.NewMemoryRepository()
	service := user.NewService(repo)
	handler := httpapi.NewHandler(service)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("users-api listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
