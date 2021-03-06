package author

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/silphid/readcommend/src/server/internal/db"
)

// Table represents the `author` database table, along with all
// operations that can be performed against it
type Table struct {
	queryer db.Queryer
}

// Author represents the information about one author
type Author struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// NewTable creates a new Table object using given queryer to access database
func NewTable(queryer db.Queryer) Table {
	return Table{queryer: queryer}
}

// GetAll retrieves all authors from author table, ordered by first then last name
func (a Table) GetAll(ctx context.Context) ([]Author, error) {
	// Building query
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := builder.
		Select("id", "first_name", "last_name").
		From("author").
		OrderBy("first_name").
		OrderBy("last_name").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %w", err)
	}

	// Execute query
	rows, err := a.queryer.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	// Scan rows and collect authors
	authors := []Author{}
	for rows.Next() {
		author := Author{}
		if err := rows.Scan(
			&author.ID,
			&author.FirstName,
			&author.LastName,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		authors = append(authors, author)
	}

	return authors, nil
}
