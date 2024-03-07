package handler

import (
	"errors"
	"math/rand"

	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/Hunter-club/cloudman/view"
	"github.com/labstack/echo/v4"
)

// PreAllocateLine 这个接口是第一个调用的
func PreAllocateLine(c echo.Context) (interface{}, error) {
	var err error

	// 获取预分配资源的请求
	req := &view.AllocateRequest{}

	err = c.Bind(req)
	if err != nil {
		return nil, err
	}
	db := database.GetDB()
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	err = AllocateForZone(req, c)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func AllocateForZone(req *view.AllocateRequest, c echo.Context) error {
	db := database.GetDB()
	hostDeviceAllocate := make([]*models.HostDeviceAllocate, 0)
	for zone, quota := range req.Lines {
		unallocatedHosts, err := getUnallocatedHosts(zone)
		if err != nil {
			return err
		}
		if len(unallocatedHosts) < quota {
			return errors.New("IP not enough")
		}
		hosts := selectRandomHosts(unallocatedHosts, quota)
		for _, host := range hosts {
			hostDeviceAllocate = append(hostDeviceAllocate, &models.HostDeviceAllocate{
				HostID:  host.HostID,
				OrderID: req.OrderID,
			})
		}
	}
	return db.Create(hostDeviceAllocate).Error
}

func getUnallocatedHosts(zone string) ([]models.Host, error) {
	db := database.GetDB()
	var unallocatedHosts []models.Host
	if err := db.Model(&models.Host{}).
		Joins("LEFT JOIN host_device_allocates ON hosts.host_id = host_device_allocates.host_id").
		Where("host_device_allocates.host_id IS NULL").
		Where("hosts.zone = ?", zone).
		Select("hosts.*").
		Find(&unallocatedHosts).Error; err != nil {
		return nil, err
	}
	return unallocatedHosts, nil
}

func selectRandomHosts(hosts []models.Host, n int) []models.Host {
	rand.Shuffle(len(hosts), func(i, j int) {
		hosts[i], hosts[j] = hosts[j], hosts[i]
	})
	if n > len(hosts) {
		return hosts // 如果n大于hosts的长度，返回整个切片
	}
	return hosts[:n] // 返回随机选取的n个hosts
}
