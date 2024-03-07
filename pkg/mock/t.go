package main

import (
	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/google/uuid"
)

func WriteHost() {
	db := database.GetDB()
	err := db.Create(&models.Host{
		Name:      "localhost",
		PrimaryIP: "127.0.0.1",
		Zone:      "lan",
		HostID:    uuid.NewString(),
	}).Error

	if err != nil {
		panic(err)
	}
}

func main() {
	WriteHost()
}
