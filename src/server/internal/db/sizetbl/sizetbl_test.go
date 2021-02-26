package sizetbl

import (
	"context"
	"os"
	"testing"

	"github.com/silphid/readcommend/src/server/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAll(t *testing.T) {
	ctx := context.Background()
	dbURL := os.Getenv("DB_URL")
	database, err := db.New(ctx, dbURL)
	require.NoError(t, err)

	table := New(database)
	sizes, err := table.GetAll(ctx)
	require.NoError(t, err)

	assert.Equal(t, 7, len(sizes), "number of results")

	first := sizes[0]
	second := sizes[1]
	third := sizes[2]
	last := sizes[len(sizes)-1]
	assert.Equal(t, "Any", first.Title, "sorting by ID: Any comes first")
	assert.Equal(t, "Short story – up to 35 pages", second.Title, "sorting by ID: Short story comes second")
	assert.Equal(t, "Novelette – 35 to 85 pages", third.Title, "sorting by ID: Novelette comes third")
	assert.Equal(t, "Monument – 800 pages and up", last.Title, "sorting by ID: Monument comes last")

	assert.Nil(t, second.MinPages, "NULL values are loaded as nil")
	assert.Equal(t, 34, *second.MaxPages, "regular values are loaded correctly")
}
