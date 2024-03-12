package server

import (
	"github.com/Hunter-club/cloudman/handler"
	"github.com/labstack/echo/v4"
)

func NewSub(e *echo.Echo) {

	e.GET("/sub/:sub_id", func(c echo.Context) error {
		handler.Sub(c)
		return nil
	})
	e.POST("/sub/delete", Handler(handler.DeleteSub))
	e.Start(":9999")
}
