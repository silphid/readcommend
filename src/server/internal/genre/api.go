package genre

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// API represents the REST API/HTTP interfacing part of the package.
type API struct {
	service Service
}

// NewAPI creates an API object to wrap given service.
func NewAPI(service Service) API {
	return API{
		service: service,
	}
}

// GetAll returns all items as a JSON list
func (api API) GetAll(c echo.Context) error {
	items, err := api.service.GetAll(c.Request().Context())
	if err != nil {
		return fmt.Errorf("querying service: %w", err)
	}
	return c.JSON(http.StatusOK, items)
}
