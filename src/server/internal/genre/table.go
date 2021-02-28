package genre

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/silphid/readcommend/src/server/internal/db"
)

// Table represents the `genre` database table, along with all
// operations that can be performed against it
type Table struct {
	queryer db.Queryer
}

// Genre represents the information about one genre
type Genre struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// NewTable creates a new GenreTable object using given queryer to access database
func NewTable(queryer db.Queryer) Table {
	return Table{queryer: queryer}
}

// GetAll retrieves all genres from genre table, ordered by ID
func (a Table) GetAll(ctx context.Context) ([]Genre, error) {
	// Building query
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := builder.
		Select("id", "title").
		From("genre").
		OrderBy("id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %w", err)
	}

	// Execute query
	rows, err := a.queryer.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	// Scan rows and collect genres
	genres := []Genre{}
	for rows.Next() {
		genre := Genre{}
		if err := rows.Scan(
			&genre.ID,
			&genre.Title,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		genres = append(genres, genre)
	}

	return genres, nil
}
