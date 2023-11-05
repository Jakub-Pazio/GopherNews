package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/", "./templates")
	e.Static("/htmx", "./htmx")

	// Later refacot to sever go templates filed with articles data
	e.GET("/articles", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h3>Articles</h3>")
	})

	e.Logger.Fatal(e.Start(":9999"))
}
