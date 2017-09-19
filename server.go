package main

import (
	"net/http"
	"html/template"
	"github.com/labstack/echo"
	"io"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"fmt"
	"os"
)

type (
	Users struct {
		Id      int64   `db:"id"`
		Name    string  `db:"name"`
	}
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

	conn, err := dbr.Open("mysql", "root@tcp(127.0.0.1:3306)/growth_development", nil)
	if err != nil{
		fmt.Fprintf(os.Stderr, "エラー：%d", err)
		os.Exit(1)
	}

	session := conn.NewSession(nil)

	var user []Users

	session.Select("*").From("users").Load(&user)

	fmt.Printf("%s", user)

	e.GET("/welcome", func(c echo.Context) error {
		// テンプレートに渡す値
		data := struct {
			ServiceInfo
			Content string
		} {
			ServiceInfo: serviceInfo,
			Content:  user[0].Name,
		}
		return c.Render(http.StatusOK, "welcome", data)
	})

	e.Logger.Fatal(e.Start(":1323"))
}