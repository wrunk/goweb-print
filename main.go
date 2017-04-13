package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Create our echo webserver
	e := echo.New()

	// Use the standard echo request/access logger for all requests
	e.Use(middleware.Logger())

	// Mandrill web server group. Will only serve endpoints under /mandrill/...
	g := e.Group("/mandrill")

	// We use a full request logging middleware to see exactly what mandrill will send
	g.Use(FullRequestLog())

	// Configure the only endpoint for mandrill group, a verify endpoint
	g.POST("/verify", verifyMandrill)

	// Also serve this on get for ease of testing with curl and browser
	g.GET("/verify", verifyMandrill)

	// Start the webserver
	e.Logger.Fatal(e.Start(":1323"))
}

func verifyMandrill(ctx echo.Context) error {
	return ctx.JSONPretty(http.StatusOK, map[string]interface{}{"success": true}, "  ")
}
