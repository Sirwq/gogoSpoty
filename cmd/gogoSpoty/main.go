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

const pollInterval = 5 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	application := app.New(ctx, *cfg)

	<-ctx.Done()
	log.Println("Shutting down...")

	log.Println("Stopped")
}
