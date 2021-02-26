package sizetbl

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
	sizes, err := table.GetAll(ctx)
	require.NoError(err)

	assert.Equal(7, len(sizes), "number of results")

	first := sizes[0]
	second := sizes[1]
	third := sizes[2]
	last := sizes[len(sizes)-1]
	assert.Equal("Any", first.Title, "sorting by ID: Any comes first")
	assert.Equal("Short story – up to 35 pages", second.Title, "sorting by ID: Short story comes second")
	assert.Equal("Novelette – 35 to 85 pages", third.Title, "sorting by ID: Novelette comes third")
	assert.Equal("Monument – 800 pages and up", last.Title, "sorting by ID: Monument comes last")

	assert.Nil(second.MinPages, "NULL values are loaded as nil")
	assert.Equal(34, *second.MaxPages, "regular values are loaded correctly")
}
