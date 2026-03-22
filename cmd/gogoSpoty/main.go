package main

import (
	"context"
	"gogoSpoty/internal/app"
	"gogoSpoty/internal/config"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	application, err := app.New(ctx, cfg)

	if err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	log.Printf("listening on %v", cfg.Spotify.RedirectURL)
	application.Start(ctx)

	<-ctx.Done()
	log.Println("Shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	application.Shutdown(shutdownCtx)

	log.Println("Stopped")
}
