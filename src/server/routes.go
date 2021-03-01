package main

import (
	"github.com/labstack/echo"
	"github.com/silphid/readcommend/src/server/internal/author"
	"github.com/silphid/readcommend/src/server/internal/book"
	"github.com/silphid/readcommend/src/server/internal/db"
	"github.com/silphid/readcommend/src/server/internal/era"
	"github.com/silphid/readcommend/src/server/internal/genre"
	"github.com/silphid/readcommend/src/server/internal/size"
)

func setupRoutes(root *echo.Group, db db.DB) {
	v1 := root.Group("/api/v1")
	author.SetupRoutes(v1, author.NewAPI(author.NewService(author.NewTable(db))))
	book.SetupRoutes(v1, book.NewAPI(book.NewService(book.NewTable(db))))
	era.SetupRoutes(v1, era.NewAPI(era.NewService(era.NewTable(db))))
	genre.SetupRoutes(v1, genre.NewAPI(genre.NewService(genre.NewTable(db))))
	size.SetupRoutes(v1, size.NewAPI(size.NewService(size.NewTable(db))))

	root.Static("/app", "app/dist")
}
