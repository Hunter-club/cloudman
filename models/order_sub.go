package models

import "gorm.io/gorm"

type OrderSub struct {
	gorm.Model
	OrderID string `gorm:"column:order_id"`
	SubID   string `gorm:"column:sub_id"`
	Vmess   string `gorm:"column:vmess"`
}
