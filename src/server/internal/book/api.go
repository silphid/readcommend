package book

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

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

type jsonError struct {
	Message string `json:"message"`
}

// GetRecommendations retrieves books ordered from best to worst
// rating and filtered by multiple optional criteria.
func (api API) GetRecommendations(c echo.Context) error {
	// Parse criteria from query params
	criteria := Criteria{}
	err := firstError(
		getInts(c, "authors", &criteria.authorIDs),
		getInts(c, "genres", &criteria.genreIDs),
		getUint(c, "min-pages", &criteria.minPages, 1, 10000),
		getUint(c, "max-pages", &criteria.maxPages, 1, 10000),
		getUint(c, "min-year", &criteria.minYear, 1800, 2100),
		getUint(c, "max-year", &criteria.maxYear, 1800, 2100),
		getUint(c, "limit", &criteria.limit, 1, math.MaxUint16),
	)
	if err != nil {
		result := &jsonError{
			Message: fmt.Sprintf("%s", err),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	// Delegate call to book service
	items, err := api.service.GetRecommendations(c.Request().Context(), criteria)
	if err != nil {
		return fmt.Errorf("querying service: %w", err)
	}
	return c.JSON(http.StatusOK, items)
}

// getUintsParam retrieves, parses and validates that given optional
// query param is a list of ints.
func getInts(c echo.Context, name string, values *[]int) error {
	str := c.QueryParam(name)
	if str == "" {
		*values = nil
		return nil
	}
	strs := strings.Split(str, ",")

	ints := make([]int, len(strs))
	for i, str := range strs {
		// Convert to uint
		ui, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("argument %q must be a list of integers: %w", name, err)
		}
		ints[i] = ui
	}

	*values = ints
	return nil
}

// getUint retrieves, parses and validates that given optional
// query param is within given range.
func getUint(c echo.Context, name string, value **uint64, min uint64, max uint64) error {
	str := c.QueryParam(name)
	if str == "" {
		*value = nil
		return nil
	}

	// Convert to uint
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fmt.Errorf("argument %q must be an integer", name)
	}

	// Validate range, if boundaries specified
	if i < min {
		return fmt.Errorf("argument %q must have a minimum value of %d", name, min)
	}
	if i > max {
		return fmt.Errorf("argument %q must have a maximum value of %d", name, max)
	}

	*value = &i
	return nil
}

// firstError returns first error in given errors, if any
func firstError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
