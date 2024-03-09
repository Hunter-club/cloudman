package handler

import (
	"net/http"

	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/labstack/echo/v4"
)

func Sub(c echo.Context) (interface{}, error) {

	subID := c.Param("sub_id")

	if subID == "" {
		c.String(http.StatusNotFound, "")
	}

	c.Response().Writer.Header().Set("Profile-Update-Interval", "12")
	c.Response().Writer.Header().Set("Profile-Title", subID)

	db := database.GetDB()

	orderSub := &models.OrderSub{}

	err := db.Model(&models.OrderSub{}).
		Where(&models.OrderSub{
			SubID: subID,
		}).
		Find(orderSub).Error
	if err != nil {
		c.String(http.StatusInternalServerError, "")
	}

	c.String(http.StatusOK, orderSub.Vmess)

	return nil, nil

}
