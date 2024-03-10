package server

import (
	"net/http"

	"github.com/Hunter-club/cloudman/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunServer() {
	// Echo instance
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			secret := c.Request().Header.Get("secret")

			if secret == "FXf4nzFzax8A.k-a" {
				return next(c)
			} else {
				return c.JSON(http.StatusUnauthorized, "")
			}
		}
	})
	e.POST("/api/v1/sub", Handler(handler.AllocateResource))
	e.POST("/api/v1/xray", Handler(handler.XUIConfigure))
	e.POST("/api/v1/line", Handler(handler.PreAllocateLine))
	go func() {
		e := echo.New()
		NewSub(e)
	}()
	e.Logger.Fatal(e.Start(":8080"))
}
