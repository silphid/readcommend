package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/silphid/readcommend/src/server/internal/db"
)

func main() {
	// Start by listening to system signals, first thing in the morning,
	// to avoid missing any, even during startup.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	cfg := getConfig()

	// Create context that gets cancelled by different system signals
	// and a channel to trigger graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	shutdownChan := make(chan struct{})
	go func() {
		sig := <-sigChan
		log.Printf("graceful shutdown initiated by signal: %s", sig.String())
		cancel()
		close(shutdownChan)
	}()

	// Create database connection, cancellable by context
	db, err := db.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to create database instance: %v", err)
	}

	// Create and start server
	server, err := newServer(cfg, *db, shutdownChan)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	if err := server.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
