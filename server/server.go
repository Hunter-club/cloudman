package server

import (
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
	e.POST("/api/v1/sub", Handler(handler.AllocateResource))
	e.POST("/api/v1/xray", Handler(handler.XUIConfigure))
	e.POST("/api/v1/line", Handler(handler.PreAllocateLine))

	go func() {
		e := echo.New()
		NewSub(e)
	}()

	e.Logger.Fatal(e.Start(":8080"))

}
