package genretbl

import (
	"context"
	"os"
	"testing"

	"github.com/silphid/readcommend/src/server/internal/db"
	_assert "github.com/stretchr/testify/assert"
	_require "github.com/stretchr/testify/require"
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
	table := New(database)

	// Execute query
	genres, err := table.GetAll(ctx)
	require.NoError(err)

	assert.Equal(8, len(genres), "number of results")

	first := genres[0]
	second := genres[1]
	last := genres[len(genres)-1]
	assert.Equal("Young Adult", first.Title, "sorting by ID: Young Adult comes first")
	assert.Equal("SciFi/Fantasy", second.Title, "sorting by ID: SciFi/Fantasy comes second")
	assert.Equal("Childrens", last.Title, "sorting by ID: Childrens comes last")
}
