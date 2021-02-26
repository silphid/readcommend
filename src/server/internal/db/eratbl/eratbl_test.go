package eratbl

import (
	"context"
	"os"
	"testing"

	"github.com/silphid/readcommend/src/server/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAll(t *testing.T) {
	req := require.New(t)
	assrt := assert.New(t)

	ctx := context.Background()
	dbURL := os.Getenv("DB_URL")
	database, err := db.New(ctx, dbURL)
	req.NoError(err)

	table := New(database)
	eras, err := table.GetAll(ctx)
	req.NoError(err)

	assrt.Equal(3, len(eras), "number of results")

	first := eras[0]
	second := eras[1]
	last := eras[len(eras)-1]
	assrt.Equal("Any", first.Title, "sorting by ID: Any comes first")
	assrt.Equal("Classic", second.Title, "sorting by ID: Classic comes second")
	assrt.Equal("Modern", last.Title, "sorting by ID: Modern comes last")

	assrt.Nil(second.MinYear, "NULL values are loaded as nil")
	assrt.Equal(1969, *second.MaxYear, "regular values are loaded correctly")
}
