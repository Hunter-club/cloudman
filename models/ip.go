package models

import "gorm.io/gorm"

// ip_id	UUID	IP唯一标识
// address	String	IP地址
// host_id	UUID	所属主机标识

type IP struct {
	gorm.Model
	IPID    string `gorm:"column:ip_id;type:uuid;primaryKey"`
	Address string `gorm:"column:address"`
	HostID  string `gorm:"column:host_id"`
}
