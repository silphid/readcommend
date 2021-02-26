package booktbl

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/silphid/readcommend/src/server/internal/db"
	"github.com/silphid/readcommend/src/server/internal/db/authortbl"
	"github.com/silphid/readcommend/src/server/internal/db/genretbl"
)

// BookTable represents the `book` database table, along with all
// operations that can be performed against it
type BookTable struct {
	queryer db.Queryer
}

// Book represents the information about one book
type Book struct {
	ID            int              `json:"id"`
	Title         string           `json:"title"`
	YearPublished int              `json:"yearPublished"`
	Rating        float32          `json:"rating"`
	Pages         int              `json:"pages"`
	Genre         genretbl.Genre   `json:"genre"`
	Author        authortbl.Author `json:"author"`
}

// New creates a new BookTable object using given queryer to access database
func New(queryer db.Queryer) BookTable {
	return BookTable{queryer: queryer}
}

// Criteria describes all possible criteria for querying book recommendations
type Criteria struct {
	authorIDs []int
	genreIDs  []int
	minPages  *int
	maxPages  *int
	minYear   *int
	maxYear   *int
	limit     *uint64
}

// GetRecommendations retrieves books from book table, ordered from best to worst
// rating and filtered by multiple optional criteria. Specifying an empty
// slice of authorIDs or genreIDs means those criteria will not be applied,
// whereas specifying multiple IDs means the union of results will be
// returned. For page, year and limit criteria, specifying a nil pointer
// will ignore the criteria.
func (a BookTable) GetRecommendations(ctx context.Context, criteria Criteria) ([]Book, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("b.id", "b.title", "b.year_published", "b.rating", "b.pages",
			"b.genre_id", "g.title", "b.author_id", "a.first_name", "a.last_name").
		From("book b").
		LeftJoin("author a ON a.id = b.author_id").
		LeftJoin("genre g ON g.id = b.genre_id").
		OrderBy("rating DESC")

	// Authors & genres
	cr := criteria
	builder = whereColumnInIDs(builder, "author_id", cr.authorIDs)
	builder = whereColumnInIDs(builder, "genre_id", cr.genreIDs)

	// Number of pages
	if cr.minPages != nil {
		builder = builder.Where("pages > ?", *cr.minPages)
	}
	if cr.maxPages != nil {
		builder = builder.Where("pages < ?", *cr.maxPages)
	}

	// Year published
	if cr.minYear != nil {
		builder = builder.Where("year_published > ?", *cr.minYear)
	}
	if cr.maxYear != nil {
		builder = builder.Where("year_published < ?", *cr.maxYear)
	}

	// Limit
	if cr.limit != nil {
		builder = builder.Limit(*cr.limit)
	}

	// Generate query string
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %w", err)
	}

	// Execute query
	rows, err := a.queryer.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	// Scan rows and collect books
	books := []Book{}
	for rows.Next() {
		book := Book{
			Genre:  genretbl.Genre{},
			Author: authortbl.Author{},
		}
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.YearPublished,
			&book.Rating,
			&book.Pages,
			&book.Genre.ID,
			&book.Genre.Title,
			&book.Author.ID,
			&book.Author.FirstName,
			&book.Author.LastName,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		books = append(books, book)
	}

	return books, nil
}

// whereColumnInIDs adds to given builder a "WHERE column IN (ids...)" clause
func whereColumnInIDs(builder sq.SelectBuilder, column string, ids []int) sq.SelectBuilder {
	if len(ids) == 0 {
		return builder
	}

	// Convert []int to []interface{}
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	return builder.Where(fmt.Sprintf("%s IN (%s)", column, sq.Placeholders(len(ids))), args...)
}
