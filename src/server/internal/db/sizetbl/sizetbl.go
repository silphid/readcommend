package sizetbl

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/silphid/readcommend/src/server/internal/db"
)

// SizeTable represents the `size` database table, along with all
// operations that can be performed against it
type SizeTable struct {
	queryer db.Queryer
}

// Size represents the information about one size
type Size struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	MinPages *int   `json:"min_pages"`
	MaxPages *int   `json:"max_pages"`
}

// New creates a new SizeTable object using given queryer to access database
func New(queryer db.Queryer) SizeTable {
	return SizeTable{queryer: queryer}
}

// GetAll retrieves all sizes from size table, ordered by ID
func (a SizeTable) GetAll(ctx context.Context) ([]Size, error) {
	// Building query
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := builder.
		Select("id", "title", "min_pages", "max_pages").
		From("size").
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

	// Scan rows and collect sizes
	sizes := []Size{}
	for rows.Next() {
		size := Size{}
		if err := rows.Scan(
			&size.ID,
			&size.Title,
			&size.MinPages,
			&size.MaxPages,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		sizes = append(sizes, size)
	}

	return sizes, nil
}
