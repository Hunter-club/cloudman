package handler

import (
	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/Hunter-club/cloudman/view"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func HostImport(c echo.Context) (interface{}, error) {
	req := &view.HostImportRequest{}
	err := c.Bind(req)
	if err != nil {
		return nil, err
	}

	hosts := make([]*models.Host, 0)

	db := database.GetDB()

	for _, host := range req.Hosts {
		hosts = append(hosts, &models.Host{
			Name:      host.Name,
			PrimaryIP: host.PrimaryIP,
			Zone:      host.Zone,
			HostID:    uuid.NewString(),
			Domain:    host.Domain,
		})
	}

	err = db.Create(hosts).Error

	if err != nil {
		return nil, err
	}
	return nil, nil
}
