// package api
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Tasks2/hmis/internal/config"
	"github.com/Tasks2/hmis/internal/database"
	"github.com/Tasks2/hmis/internal/logger"
	"github.com/Tasks2/hmis/internal/middleware"
	"github.com/Tasks2/hmis/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	appLogger := logger.New()

	err := godotenv.Load()
	if err != nil {
		appLogger.Info.Println("No .env file found; using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	defer db.Close()

	appRouter := router.New(db)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.RequestLogger(appRouter),
	}

	log.Printf("HMIS API listening on http://localhost:%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
