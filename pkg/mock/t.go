package main

import (
	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/google/uuid"
)

func WriteHost() {

	hosts := make([]*models.Host, 0)

	hosts = append(hosts, &models.Host{
		Name:      "gw.gcore.lu",
		PrimaryIP: "92.38.149.1",
		Zone:      "us",
		HostID:    uuid.NewString(),
	})

	hosts = append(hosts, &models.Host{
		Name:      "10.255.52.183",
		Zone:      "us",
		PrimaryIP: "10.255.52.183",
		HostID:    uuid.NewString(),
	})

	hosts = append(hosts, &models.Host{
		Name:      "10.255.52.188",
		Zone:      "us",
		PrimaryIP: "10.255.52.188",
		HostID:    uuid.NewString(),
	})

	hosts = append(hosts, &models.Host{
		Name:      "palo-b24-link.ip.twelve99.net",
		PrimaryIP: "62.115.182.214	",
		Zone:      "us",
		HostID:    uuid.NewString(),
	})

	db := database.GetDB()
	err := db.Create(hosts).Error
	if err != nil {
		panic(err)
	}
}

func main() {
	WriteHost()
}
