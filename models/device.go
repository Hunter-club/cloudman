package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model            // Gorm的Model包含基本的ID, CreatedAt等字段
	DeviceID   string     `gorm:"column:device_id;type:uuid;primaryKey"` // 设备唯一标识，假设为UUID格式
	DeviceName string     `gorm:"column:device_name"`                    // 设备名称
	DeviceSpec DeviceSpec `gorm:"column:device_spec"`                    // 设备规格说明
	IsDeleted  bool       `gorm:"column:is_deleted"`                     // 是否删除
}

type DeviceSpec struct {
	Zone         string `gorm:"column:zone"`
	IPQuota      int    `gorm:"column:ip_quota"`      // IP支持个数，假设为整型
	TrafficQuota int    `gorm:"column:traffic_quota"` // 它能分配的流量
}
