package era

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
	table := NewTable(database)

	// Execute query
	eras, err := table.GetAll(ctx)
	require.NoError(err)

	assert.Equal(3, len(eras), "number of results")

	first := eras[0]
	second := eras[1]
	last := eras[len(eras)-1]
	assert.Equal("Any", first.Title, "sorting by ID: Any comes first")
	assert.Equal("Classic", second.Title, "sorting by ID: Classic comes second")
	assert.Equal("Modern", last.Title, "sorting by ID: Modern comes last")

	assert.Nil(second.MinYear, "NULL values are loaded as nil")
	assert.Equal(1969, *second.MaxYear, "regular values are loaded correctly")
}
