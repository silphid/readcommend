package era

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/silphid/readcommend/src/server/internal/db"
)

// Table represents the `era` database table, along with all
// operations that can be performed against it
type Table struct {
	queryer db.Queryer
}

// Era represents the information about one era
type Era struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	MinYear *int   `json:"min_year"`
	MaxYear *int   `json:"max_year"`
}

// NewTable creates a new Table object using given queryer to access database
func NewTable(queryer db.Queryer) Table {
	return Table{queryer: queryer}
}

// GetAll retrieves all eras from era table, ordered by ID
func (a Table) GetAll(ctx context.Context) ([]Era, error) {
	// Building query
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := builder.
		Select("id", "title", "min_year", "max_year").
		From("era").
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

	// Scan rows and collect eras
	eras := []Era{}
	for rows.Next() {
		era := Era{}
		if err := rows.Scan(
			&era.ID,
			&era.Title,
			&era.MinYear,
			&era.MaxYear,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		eras = append(eras, era)
	}

	return eras, nil
}
