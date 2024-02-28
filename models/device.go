package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model        // Gorm的Model包含基本的ID, CreatedAt等字段
	DeviceID   string `gorm:"column:device_id;type:uuid;primaryKey"` // 设备唯一标识，假设为UUID格式
	DeviceName string `gorm:"column:device_name"`                    // 设备名称
	DeviceSpec string `gorm:"column:device_spec"`                    // 设备规格说明
	IPQuota    int    `gorm:"column:ip_quota"`                       // IP支持类型，假设为整型
}
