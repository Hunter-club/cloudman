package models

import "gorm.io/gorm"

type Transfer struct {
	gorm.Model
	TransferID string `gorm:"column:transfer_id"`
	Addr       string `gorm:"column:addr"`
	Port       int    `gorm:"column:port"`
	HostID     string `gorm:"column:host_id"`
	HostIP     string `gorm:"column:host_ip"`
	VmessID    string `gorm:"column:vmess_id"`
	OrderID    string `gorm:"column:order_id"`
}
