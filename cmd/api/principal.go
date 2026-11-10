package main

import (
	"log"
	"minimarket/internal/application"
	"minimarket/internal/infrastructure"
	"minimarket/internal/presentation"
	"net/http"
	"os"
	"time"
)

func main() {
	var repo application.Repository
	if os.Getenv("STORAGE_ADAPTER") == "memory" {
		memRepo, err := infrastructure.NewMemoryAdapter()
		if err != nil {
			log.Fatal(err)
		}
		repo = memRepo
	} else {
		path := os.Getenv("DATA_FILE")
		if path == "" {
			path = "data.json"
		}
		jsonRepo, err := infrastructure.NewJSONRepository(path)
		if err != nil {
			log.Fatal(err)
		}
		if err := jsonRepo.EnsureRichSeed(); err != nil {
			log.Printf("Aviso al inicializar semilla: %v", err)
		}
		repo = jsonRepo
	}

	handler := presentation.NewHandler(application.Service{Repo: repo})

	// En servidores Linux EC2 sin Caddy activo, también atiende directamente en el puerto 80 si está disponible
	if os.Getenv("PUBLIC_HTTP_80") == "1" {
		go func() {
			pub := &http.Server{
				Addr:              ":80",
				Handler:           handler,
				ReadHeaderTimeout: 5 * time.Second,
				ReadTimeout:       10 * time.Second,
				WriteTimeout:      10 * time.Second,
				IdleTimeout:       60 * time.Second,
			}
			if err := pub.ListenAndServe(); err != nil {
				log.Printf("Puerto 80 gestionado externamente (ej. Caddy): %v", err)
			}
		}()
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("Minimarket SaaS Multiempresa disponible en http://%s", addr)
	log.Fatal(server.ListenAndServe())
}
