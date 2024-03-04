package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := database.Db

	e.POST("/api/v1/device/sub", func(c echo.Context) error {
		// 获取订阅的连接
		// 为设备创建固定的规则
		// 入站、出站、路由规则
		// 生成订阅
		// todo 明天写了 今天脑子不清楚了
		return nil
	})

	e.POST("/api/v1/device/allocate", func(c echo.Context) error {
		device := &models.Device{}
		err := c.Bind(device)
		if err != nil {
			DealResp(c, "bind error", err, nil)
		}

		zone := device.DeviceSpec.Zone
		quota := device.DeviceSpec.IPQuota
		unallocatedHosts, err := getUnallocatedHosts(db, zone)
		// 1000是IP不够的预警
		if err != nil || len(unallocatedHosts) <= quota+1000 {
			DealResp(c, "get unallocated hosts error", err, nil)
		}
		hosts := selectRandomHosts(unallocatedHosts, quota)
		hostDeviceAllocate := make([]*models.HostDeviceAllocate, 0)
		for _, host := range hosts {
			hostDeviceAllocate = append(hostDeviceAllocate, &models.HostDeviceAllocate{
				HostID:   host.ID,
				DeviceID: device.ID,
			})
		}
		err = db.Create(hostDeviceAllocate).Error

		if err != nil {
			DealResp(c, "create host device allocate error", err, nil)
		}

		return c.JSON(200, JsonObj{
			Data:    nil,
			Success: true,
			Msg:     "",
		})
	})

	// Routes
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
func getUnallocatedHosts(db *gorm.DB, zone string) ([]models.Host, error) {
	var unallocatedHosts []models.Host
	if err := db.Model(&models.Host{}).
		Joins("LEFT JOIN host_device_allocates ON hosts.id = host_device_allocates.host_id").
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

type JsonObj struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
}

func DealResp(c echo.Context, msg string, err error, data interface{}) {
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, JsonObj{
			Msg:     msg,
			Success: false,
			Data:    data,
		})
	} else {
		c.JSON(http.StatusOK, JsonObj{
			Msg:     msg,
			Success: true,
			Data:    data,
		})
	}
}
