package main

import (
	"line-go-healthcheck/handler"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))
	e.POST("/upload", handler.Upload)

	e.Logger.Fatal(e.Start(":8080"))
}
