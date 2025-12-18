package main

import (
	"context"
	"log"
	"notes-service/config"
	"notes-service/database/handlers"
	"notes-service/database/repository"
	"notes-service/database/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {

	cfg := config.LoadConfig()

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	userRepo := repository.NewNoteRepository(db)
	userService := service.NewNoteService(userRepo)
	userHandler := handlers.NewNoteHandler(userService)

	app := fiber.New(fiber.Config{
		AppName:      "Notes service API",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	userHandler.RegisterRoutes(app)

	port := cfg.APPort
	if port == "" {
		port = "8081"
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Println("Shutting down server gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
	}()

	log.Printf("Server starting on http://localhost:%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("Server stopped")
}
