package booktbl

import (
	"context"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/silphid/readcommend/src/server/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testCase defines a sub-test-case to run
type testCase struct {
	name     string
	criteria Criteria

	// We can either check only the book IDs (to check presence and order of results)
	// or actual book objects (to check proper contents of book objects with deep equality)
	expectedBookIDs *[]int
	expectedBooks   *[]Book
}

// testCases defines the data for all sub-test-cases to run
var testCases = []testCase{
	{
		name: "querying multiple authors returns union of their books by rating",
		criteria: Criteria{
			authorIDs: []int{1, 2, 30},
		},
		expectedBookIDs: &[]int{27, 52, 55, 24, 29, 17, 50, 37, 10},
	},
}

func TestGetRecommendations(t *testing.T) {
	// Connect to database and get table object
	ctx := context.Background()
	dbURL := os.Getenv("DB_URL")
	database, err := db.New(ctx, dbURL)
	require.NoError(t, err)
	table := New(database)

	// Run all sub-tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Execute query
			actualBooks, err := table.GetRecommendations(ctx, tc.criteria)
			assert.NoError(t, err)

			// Check returned book IDs or the actual book objects
			if tc.expectedBookIDs != nil {
				// Test only book IDs
				actualBookIDs := getBookIDs(actualBooks)
				if diff := deep.Equal(*tc.expectedBookIDs, actualBookIDs); diff != nil {
					t.Error(diff)
				}
			} else if tc.expectedBooks != nil {
				// Test for book objects deep equality
				if diff := deep.Equal(*tc.expectedBooks, actualBooks); diff != nil {
					t.Error(diff)
				}
			} else {
				t.Fatal("must specify either expected book IDs or actual book objects")
			}
		})
	}
}

func getBookIDs(books []Book) []int {
	ids := make([]int, len(books))
	for i, book := range books {
		ids[i] = book.ID
	}
	return ids
}
