package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/Hunter-club/cloudman/view"
	"github.com/labstack/echo/v4"
)

func DeleteSub(c echo.Context) (interface{}, error) {
	req := view.SubRequest{}

	err := c.Bind(c)
	if err != nil {
		return nil, err
	}

	db := database.GetDB()

	err = db.Delete(&models.OrderSub{}, &models.OrderSub{
		OrderID: req.OrderID,
	}).Error

	if err != nil {
		return nil, err
	}

	return "deleted", nil

}

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

	vmess := base64.StdEncoding.EncodeToString([]byte(orderSub.Vmess))
	c.String(http.StatusOK, vmess)

	return nil, nil

}
