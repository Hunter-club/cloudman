package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	go startAllocateJob()
	go startOrderExpireJob()

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/api/v1/order", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	e.GET("/api/v1/device", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	e.POST("/api/v1/hosts", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// 启动一个分配协程负责订单的分配
func startAllocateJob() {
	for {
	}
}

func startOrderExpireJob() {
	for {
	}
}
