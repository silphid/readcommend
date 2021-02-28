package genre

import "github.com/labstack/echo"

// SetupRoutes maps the different HTTP routes to API implementations
func SetupRoutes(e *echo.Group, api API) {
	e.GET("/genres", api.GetAll)
}
