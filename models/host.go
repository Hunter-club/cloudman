package models

import "gorm.io/gorm"

type Host struct {
	gorm.Model
	HostID    string `gorm:"column:host_id;type:uuid;primaryKey"`
	Name      string `grom:"column:name"`
	PrimaryIP string `gorm:"column:primary_ip"`
	Zone      string `gorm:"column:zone"`
}
