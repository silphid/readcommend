package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/silphid/readcommend/src/server/internal/db"
)

// server represents an echo server that is ready to be started
type server struct {
	echo         *echo.Echo
	cfg          config
	shutdownChan <-chan struct{}
}

// newServer creates an echo server that is ready to be started
func newServer(
	cfg config,
	db db.DB,
	shutdownChan <-chan struct{},
) (*server, error) {
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

	return &server{echo: e, cfg: cfg, shutdownChan: shutdownChan}, nil
}

// Start kicks off echo server and handles graceful shutdown triggered by channel
func (s *server) Start() error {
	// Start server
	go func() {
		if err := s.echo.Start(fmt.Sprintf(":%d", s.cfg.Port)); err != nil && err != http.ErrServerClosed {
			s.echo.Logger.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for shutdown event to gracefully shutdown server, with a timeout after grace period
	<-s.shutdownChan
	s.echo.Logger.Infof("trying to gracefully shut server down within %s", s.cfg.GracePeriod.String())
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.GracePeriod)
	defer cancel()
	if err := s.echo.Shutdown(shutdownCtx); err != nil {
		s.echo.Logger.Fatalf("failed to gracefully shut server down: %v", err)
	}
	return nil
}
