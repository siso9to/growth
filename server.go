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


	e.GET("/welcome", func(c echo.Context) error {
		var user []Users
		session := conn.NewSession(nil)
		session.Select("*").From("users").Load(&user)

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

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/users/:id", func(c echo.Context) error {
		var user []Users

		session := conn.NewSession(nil)

		id := c.Param("id")

		fmt.Println(id)

		session.Select("*").
		From("users").
		Where("id = ?", id).
		Load(&user)

		fmt.Printf("%v", user)

		return c.String(http.StatusOK, user[0].Name)
	})

	e.Logger.Fatal(e.Start(":1323"))
}