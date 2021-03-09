package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/silphid/readcommend/src/server/internal/db"
)

// server represents an echo server that is ready to be started
type server struct {
	echo *echo.Echo
}

// newServer creates an echo server that is ready to be started
func newServer(db db.DB) (*server, error) {
	// Create echo server with graceful shutdown handler
	e := echo.New()
	e.Server.RegisterOnShutdown(func() {
		e.Logger.Info("closing DB connection")
		db.Close()
	})

	// Register middlewares
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	// Setup routes to different APIs
	root := e.Group("")
	setupRoutes(root, db)

	return &server{e}, nil
}

// Start kicks off echo server and handles graceful shutdown triggered by channel
func (s *server) Start(port uint16) {
	go func() {
		if err := s.echo.Start(fmt.Sprintf(":%d", port)); err != nil {
			log.Printf("server exited with error: %v", err)
		}
	}()
}

// Shutdown tries to stop the server gracefully with given grace period
func (s *server) Shutdown(gracePeriod time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()
	return s.echo.Shutdown(ctx)
}
