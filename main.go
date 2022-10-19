package main

import (
	"fmt"
	"github.com/hosseintrz/hit-tracker/internal/database"
	"github.com/hosseintrz/hit-tracker/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := database.InitStore()
	if err != nil {
		log.Fatalf("couldn't open database -> \n %s \n", err.Error())
	}
	defer db.Close()

	e.GET("/ping", func(c echo.Context) error {
		response := struct {
			Status string
		}{"ok"}
		return c.JSON(200, response)
	})
	e.GET("/*", func(c echo.Context) error {
		return handler.RootHandler(db, c)
	})

	defaultPort := "9090"
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = defaultPort
	}
	fmt.Printf("listening on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
