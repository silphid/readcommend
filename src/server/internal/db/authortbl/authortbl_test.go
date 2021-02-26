package authortbl

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
	authors, err := table.GetAll(ctx)
	require.NoError(t, err)

	assert.Equal(t, 41, len(authors), "number of results")

	first := authors[0]
	second := authors[1]
	last := authors[len(authors)-1]
	assert.Equal(t, "Abraham", first.FirstName, "firstly sorting by first name: Abraham comes first")
	assert.Equal(t, "Amelia", second.FirstName, "firstly sorting by first name: Amelia comes second")
	assert.Equal(t, "Wendell", last.FirstName, "firstly sorting by first name: Wendell comes last")

	robertMilofskyIndex := getAuthorIndexByID(authors, 38)
	robertPlimptonIndex := getAuthorIndexByID(authors, 37)
	assert.Less(t, robertMilofskyIndex, robertPlimptonIndex, "secondly sorting by last name: Robert Milofsky comes before Robert Plimpton")
}

func getAuthorIndexByID(authors []Author, id int) int {
	for i, author := range authors {
		if author.ID == id {
			return i
		}
	}
	return -1
}
