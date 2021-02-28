package size

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

// GetAll retrieves all authors ordered by first then last name
func (s Service) GetAll(ctx context.Context) ([]Size, error) {
	return s.table.GetAll(ctx)
}
