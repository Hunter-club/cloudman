package models

import "gorm.io/gorm"

type Transfer struct {
	gorm.Model
	TransferID string `gorm:"column:transfer_id"`
	Addr       string `gorm:"column:addr"`
	Port       int    `gorm:"column:port"`
}
