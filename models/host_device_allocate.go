package models

import "gorm.io/gorm"

// allocate_id	UUID	分配记录唯一标识
// ip_id	String	网卡ID
// port	String	分配的端口
// device_id	UUID	设备ID (可为空)
// order_id	UUID	订单ID (可为空)
// status	Bool	分配状态 (0未分配, 1已分配)
// vmess	string	vmess串

// 比方说我现在有一个路由器，然后这个路由器,1000MB的traffic，然后它可以连5个不同区域的IP
// 然后我有十台机器分别在美国、韩国、日本，现在我要给这个路由器分配5个IP，每个机器只有一个IP，我应该怎么做？

// 我先找到哪些机器可以分配的，insert into 然后机器的id 它是unique的 只要insert失败的话 就要重试
type HostDeviceAllocate struct {
	gorm.Model
	HostID  uint   `gorm:"column:host_id;unique"` // 设备唯一标识，假设为UUID格式
	OrderID string `gorm:"column:device_id"`
}
