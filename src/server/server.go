package main

import (
	"context"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/silphid/readcommend/src/server/internal/author"
	"github.com/silphid/readcommend/src/server/internal/db"
)

// server represents an echo server that is ready to be started
type server struct {
	echo *echo.Echo
	port uint16
}

// Start kicks off echo server
func (s *server) Start() error {
	if err := s.echo.Start(fmt.Sprintf(":%d", s.port)); err != nil {
		return err
	}
	return nil
}

// newServer creates an echo server that is ready to be started
func newServer(
	ctx context.Context,
	port uint16,
	db *db.DB,
	shutdown func(),
) (*server, error) {
	// Create echo server with graceful shutdown
	e := echo.New()
	e.Server.RegisterOnShutdown(shutdown)

	// Register middlewares
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	// Setup routes to different APIs
	v1 := e.Group("/api/v1")
	author.SetupRoutes(v1, author.NewAPI(author.NewService(author.NewTable(db))))

	return &server{echo: e, port: port}, nil
}
