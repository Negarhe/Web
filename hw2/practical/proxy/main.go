package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	handler, err := newHandler(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              cfg.listenAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("jwt proxy listening on %s", cfg.listenAddr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

type config struct {
	listenAddr string
	upstream   string
	secret     []byte
	username   string
	password   string
	tokenTTL   time.Duration
}

func loadConfig() (config, error) {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 16 {
		return config{}, errors.New("JWT_SECRET must be set and at least 16 characters")
	}
	password := os.Getenv("AUTH_PASSWORD")
	if password == "" {
		return config{}, errors.New("AUTH_PASSWORD must be set")
	}

	return config{
		listenAddr: getenv("LISTEN_ADDR", ":8080"),
		upstream:   getenv("UPSTREAM_URL", "http://legacy:8081"),
		secret:     []byte(secret),
		username:   getenv("AUTH_USERNAME", "student"),
		password:   password,
		tokenTTL:   time.Hour,
	}, nil
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
