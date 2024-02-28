package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	ID         string `gorm:"primaryKey"` // 订单唯一标识 (UUID)
	OrderID    string // 商城内订单ID
	UserID     string // 用户ID
	ExpireTime int64  // 订单过期时间
	Type       int    // 续费、新订单
	Status     int    // 已处理还是未处理
}
