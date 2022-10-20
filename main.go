package main

import (
	"fmt"
	"github.com/hosseintrz/hit-tracker/internal/database"
	"github.com/hosseintrz/hit-tracker/internal/handler"
	middleware2 "github.com/hosseintrz/hit-tracker/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {

	e := echo.New()
	e.Use(middleware.Recover())
	logger := logrus.New()

	//database
	db, err := database.InitStore()
	if err != nil {
		logger.Fatalf("couldn't open database -> \n %s \n", err.Error())
	}
	defer db.Close()

	//middlewares
	md := middleware2.NewMiddleWare(logger)
	hitTrackerMD := func(next echo.HandlerFunc) echo.HandlerFunc {
		return md.HitTracker(next, db)
	}
	e.Use(hitTrackerMD)

	//handlers
	h := handler.NewHandler(logger)

	e.GET("/ping", func(c echo.Context) error {
		response := struct {
			Status string
		}{"ok"}
		return c.JSON(200, response)
	})
	e.GET("/state", func(c echo.Context) error {
		return h.GetStats(c, db)
	})
	e.GET("/login", h.LoginHandler)
	e.GET("/payment", h.PaymentHandler)
	e.GET("/*", func(c echo.Context) error {
		return c.String(200, fmt.Sprint("you hit ", c.Request().URL.Path))
	})

	defaultPort := "9090"
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = defaultPort
	}
	fmt.Printf("listening on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
