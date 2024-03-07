package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Handler(handle func(c echo.Context) (interface{}, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		data, err := handle(c)
		if err != nil {
			Resp(c, "server internal error", err, data)
		}
		return Resp(c, "successly handle", nil, data)
	}
}

type RespObj struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
}

func Resp(c echo.Context, msg string, err error, data interface{}) error {
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusOK, RespObj{
			Msg:     msg,
			Success: false,
			Data:    data,
		})
	} else {
		return c.JSON(http.StatusOK, RespObj{
			Msg:     msg,
			Success: true,
			Data:    data,
		})
	}
}
