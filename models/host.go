package models

import "gorm.io/gorm"

type Host struct {
	gorm.Model
	Name      string `grom:"column:name"`
	PrimaryIP string `gorm:"column:primary_ip"`
	Zone      string `gorm:"column:zone"`
	HostID    string `gorm:"column:host_id;unique"`
	Domain    string `gorm:"column:domain"`
}
