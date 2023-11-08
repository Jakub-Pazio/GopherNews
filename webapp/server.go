package main

import (
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
	e.GET("/articles", articlesList)
	e.POST("/articles", ArticlesKeywordList)
	e.Logger.Fatal(e.Start(":9999"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func articlesList(c echo.Context) error {
	fetchedArticles := articles.GetArticlesFromDatabase()
	return c.Render(200, "articles", fetchedArticles)
}

func ArticlesKeywordList(c echo.Context) error {
	keyword := c.FormValue("keyword")
	fetchedArticles := articles.GetArticlesFromDatabaseKeyword(keyword)
	return c.Render(200, "articles", fetchedArticles)
}
