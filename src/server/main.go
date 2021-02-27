package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/silphid/readcommend/src/server/internal/db"
	"go.uber.org/zap"
)

func main() {
	cfg := getConfig()

	// Create context that gets cancelled by different system signals
	ctx, dispose := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer dispose()

	// Create database connection
	db, err := db.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatal("failed to create database instance", zap.Error(err))
	}

	// Create and start server with graceful shutdown handler
	shutdown := func() {
		db.Close()
	}
	server, err := newServer(ctx, cfg.Port, db, shutdown)
	if err != nil {
		log.Fatal("failed to create server", zap.Error(err))
	}
	if err := server.Start(); err != nil {
		log.Fatal("failed to start server", zap.Error(err))
	}
}
