package models

import "gorm.io/gorm"

type HostTransfer struct {
	gorm.Model
	TransferID string `gorm:"column:transfer_id;unique"`
	HostID     string `gorm:"column:host_id;unique"` // 设备唯一标识，假设为UUID格式
}
