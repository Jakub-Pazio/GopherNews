package main

import (
	// "net/http"
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopher_news/internal/articles"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/", "./templates")
	e.Static("/htmx", "./htmx")

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	// Later refactor to sever go templates filed with articles data
	e.Renderer = t
	e.GET("/articles", ArticlesList)
	e.Logger.Fatal(e.Start(":9999"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func ArticlesList(c echo.Context) error {
	// Log that finction was called, TODO: CHANGE THIS
	fmt.Println("ArticlesList called")
	// TODO: Fetch articles from database
	fetchedArticles := []articles.Article{
		articles.NewArticle("Learn C++", "https://learn-cpp.org", "C++"),
		articles.NewArticle("Learn Go", "https://learn-go.org", "Go"),
	}
	return c.Render(200, "articles", fetchedArticles)
}
