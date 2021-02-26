package booktbl

import (
	"context"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/silphid/readcommend/src/server/internal/db"
	"github.com/silphid/readcommend/src/server/internal/db/authortbl"
	"github.com/silphid/readcommend/src/server/internal/db/genretbl"
	_assert "github.com/stretchr/testify/assert"
	_require "github.com/stretchr/testify/require"
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

// pInt is a helper to convert an `int` literal into a `*int`.
// This is necessary because we use `*int` as nullable search criteria.
func pInt(value int) *int {
	return &value
}

// pInt is a helper to convert an `uint64` literal into a `*uint64`.
// This is necessary because we use `*uint64` as nullable limit criteria.
func pUint64(value uint64) *uint64 {
	return &value
}

// testCases defines the data for all sub-test-cases to run
var testCases = []testCase{
	{
		name:     "no criteria",
		criteria: Criteria{},
		expectedBookIDs: &[]int{12, 32, 27, 57, 31, 52, 56, 53, 8, 55, 9, 23, 16,
			44, 24, 21, 28, 45, 15, 29, 13, 39, 18, 7, 22, 19, 20, 4, 40, 38, 30,
			46, 14, 48, 42, 6, 17, 43, 35, 50, 54, 2, 34, 41, 49, 5, 37, 26, 25, 33,
			1, 36, 47, 10, 3, 11, 58, 51},
	},
	{
		name: "multiple authors",
		criteria: Criteria{
			authorIDs: []int{1, 2, 30},
		},
		expectedBookIDs: &[]int{27, 52, 55, 24, 29, 17, 50, 37, 10},
	},
	{
		name: "multiple authors and single genre",
		criteria: Criteria{
			authorIDs: []int{1, 2, 30},
			genreIDs:  []int{3},
		},
		expectedBookIDs: &[]int{27, 52, 55, 24, 50, 37, 10},
	},
	{
		name: "single author",
		criteria: Criteria{
			authorIDs: []int{10},
		},
		expectedBookIDs: &[]int{32, 46},
	},
	{
		name: "multiple genres",
		criteria: Criteria{
			genreIDs: []int{5, 6},
		},
		expectedBookIDs: &[]int{32, 57, 56, 21, 46, 58},
	},
	{
		name: "single genre",
		criteria: Criteria{
			genreIDs: []int{5},
		},
		expectedBookIDs: &[]int{32, 46, 58},
	},
	{
		name: "min pages",
		criteria: Criteria{
			minPages: pInt(867),
		},
		expectedBookIDs: &[]int{31, 56, 17},
	},
	{
		name: "max pages",
		criteria: Criteria{
			maxPages: pInt(16),
		},
		expectedBookIDs: &[]int{18, 38, 54},
	},
	{
		name: "min pages is inclusive",
		criteria: Criteria{
			authorIDs: []int{10},
			minPages:  pInt(731),
		},
		expectedBookIDs: &[]int{46},
	},
	{
		name: "max pages is inclusive",
		criteria: Criteria{
			authorIDs: []int{10},
			maxPages:  pInt(449),
		},
		expectedBookIDs: &[]int{32},
	},
	{
		name: "min year",
		criteria: Criteria{
			minYear: pInt(2020),
		},
		expectedBookIDs: &[]int{27, 34},
	},
	{
		name: "max year",
		criteria: Criteria{
			maxYear: pInt(1931),
		},
		expectedBookIDs: &[]int{15, 36},
	},
	{
		name: "min year is inclusive",
		criteria: Criteria{
			authorIDs: []int{10},
			minYear:   pInt(2004),
		},
		expectedBookIDs: &[]int{46},
	},
	{
		name: "max year is inclusive",
		criteria: Criteria{
			authorIDs: []int{10},
			maxYear:   pInt(1935),
		},
		expectedBookIDs: &[]int{32},
	},
	{
		name: "limit",
		criteria: Criteria{
			authorIDs: []int{10},
			limit:     pUint64(1),
		},
		expectedBookIDs: &[]int{32},
	},
	{
		name: "all book metadata properly populated",
		criteria: Criteria{
			authorIDs: []int{25, 26},
		},
		expectedBooks: &[]Book{
			{
				ID:            53,
				Title:         "Turn Left Til You Get There",
				YearPublished: 1985,
				Rating:        4.54,
				Pages:         331,
				Genre: genretbl.Genre{
					ID:    7,
					Title: "Fiction",
				},
				Author: authortbl.Author{
					ID:        26,
					FirstName: "Kris",
					LastName:  "Elegant",
				},
			},
			{
				ID:            3,
				Title:         "A Horrible Human with the Habits of a Monster",
				YearPublished: 1976,
				Rating:        1.14,
				Pages:         258,
				Genre: genretbl.Genre{
					ID:    7,
					Title: "Fiction",
				},
				Author: authortbl.Author{
					ID:        25,
					FirstName: "Kenneth",
					LastName:  "Douglas",
				},
			},
		},
	},
}

func TestGetRecommendations(t *testing.T) {
	// Helpers
	require := _require.New(t)

	// Connect to database and get table object
	ctx := context.Background()
	dbURL := os.Getenv("DB_URL")
	database, err := db.New(ctx, dbURL)
	require.NoError(err)
	table := New(database)

	// Run all sub-tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Helpers
			assert := _assert.New(t)

			// Execute query
			actualBooks, err := table.GetRecommendations(ctx, tc.criteria)
			assert.NoError(err)

			// Check returned book IDs or the actual book objects
			if tc.expectedBookIDs != nil {
				// Test only book IDs
				actualBookIDs := getBookIDs(actualBooks)
				assert.Equal(*tc.expectedBookIDs, actualBookIDs)
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
