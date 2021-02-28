package book

import "context"

// Service represents the business-logic portion of the package, which is
// abstracted away from HTTP and DB concerns. In the case of this very
// simple app, much of the complexity resides in database access, so
// this layer is almost empty, but it finds its true value as complexity
// grows.
type Service struct {
	table Table
}

// NewService creates a service object wrapping given table.
func NewService(table Table) Service {
	return Service{
		table: table,
	}
}

// GetRecommendations retrieves books ordered from best to worst
// rating and filtered by multiple optional criteria.
func (s Service) GetRecommendations(ctx context.Context, criteria Criteria) ([]Book, error) {
	return s.table.GetRecommendations(ctx, criteria)
}
