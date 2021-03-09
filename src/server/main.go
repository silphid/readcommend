package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/silphid/readcommend/src/server/internal/db"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := getConfig()

	// Create database connection, cancellable by context
	db, err := db.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to create database instance: %v", err)
	}

	// Create and start server
	server, err := newServer(*db)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	server.Start(cfg.Port)

	// Wait for termination signal and shutdown
	<-ctx.Done()
	log.Println("termination signal received, trying to shutdown gracefully...")
	if err := server.Shutdown(cfg.GracePeriod); err != nil {
		log.Fatalf("failed to gracefully shutdown server: %v", err)
	}
	log.Println("successfully shutdown")
}
