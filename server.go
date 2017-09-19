package main

import (
	"net/http"
	"html/template"
	"github.com/labstack/echo"
	"io"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// サイトで共通情報
type ServiceInfo struct {
	Title string
}

var serviceInfo = ServiceInfo {
	"Growth - Simple workflow tool for 1on1",
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e := echo.New()

	e.Renderer = t

	e.GET("/welcome", func(c echo.Context) error {
		// テンプレートに渡す値
		data := struct {
			ServiceInfo
			Content string
		} {
			ServiceInfo: serviceInfo,
			Content: "ようこそ",
		}
		return c.Render(http.StatusOK, "welcome", data)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
