package genretbl

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
	genres, err := table.GetAll(ctx)
	require.NoError(t, err)

	assert.Equal(t, 8, len(genres), "number of results")

	first := genres[0]
	second := genres[1]
	last := genres[len(genres)-1]
	assert.Equal(t, "Young Adult", first.Title, "sorting by ID: Young Adult comes first")
	assert.Equal(t, "SciFi/Fantasy", second.Title, "sorting by ID: SciFi/Fantasy comes second")
	assert.Equal(t, "Childrens", last.Title, "sorting by ID: Childrens comes last")
}
