package author

import (
	"context"
	"os"
	"testing"

	"github.com/silphid/readcommend/src/server/internal/db"
	_require "github.com/stretchr/testify/assert"
	_assert "github.com/stretchr/testify/require"
)

func TestGetAll(t *testing.T) {
	// Helpers
	require := _require.New(t)
	assert := _assert.New(t)

	// Connect to database and get table object
	ctx := context.Background()
	dbURL := os.Getenv("DB_URL")
	database, err := db.New(ctx, dbURL)
	require.NoError(err)
	table := NewTable(database)

	// Execute query
	authors, err := table.GetAll(ctx)
	require.NoError(err)

	assert.Equal(41, len(authors), "number of results")

	first := authors[0]
	second := authors[1]
	last := authors[len(authors)-1]
	assert.Equal("Abraham", first.FirstName, "firstly sorting by first name: Abraham comes first")
	assert.Equal("Amelia", second.FirstName, "firstly sorting by first name: Amelia comes second")
	assert.Equal("Wendell", last.FirstName, "firstly sorting by first name: Wendell comes last")

	robertMilofskyIndex := getAuthorIndexByID(authors, 38)
	robertPlimptonIndex := getAuthorIndexByID(authors, 37)
	assert.Less(robertMilofskyIndex, robertPlimptonIndex, "secondly sorting by last name: Robert Milofsky comes before Robert Plimpton")
}

func getAuthorIndexByID(authors []Author, id int) int {
	for i, author := range authors {
		if author.ID == id {
			return i
		}
	}
	return -1
}
