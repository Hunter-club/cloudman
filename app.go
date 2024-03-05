package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"

	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/Hunter-club/cloudman/view"
	"github.com/Hunter-club/cloudman/xui"
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
		var err error

		req := &view.SubRequest{}

		err = c.Bind(req)

		if err != nil {
			DealResp(c, "bind error", err, nil)
		}

		tx := db.Begin()
		defer func() {
			if err == nil {
				tx.Commit()
			} else {
				tx.Rollback()
			}
		}()

		ipList, err := SeletIPByOrderID(db, req.OrderID)

		if err != nil {
			return err
		}

		for _, ip := range ipList {

		}

		// todo 明天写了 今天脑子不清楚了
		return nil
	})

	e.POST("/api/v1/line/allocate", func(c echo.Context) error {
		var err error
		req := &view.AllocateRequest{}

		err = c.Bind(req)
		if err != nil {
			DealResp(c, "bind error", err, nil)
		}

		tx := db.Begin()
		defer func() {
			if err == nil {
				tx.Commit()
			} else {
				tx.Rollback()
			}
		}()
		err = AllocateForZone(db, req, c)
		if err != nil {
			DealResp(c, "allocate error", err, nil)
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

func ConfigXraySingle(host *models.Host, remark string, isDomain, isTls bool) error {
	var err error
	ip := host.PrimaryIP
	zone := host.Zone

	if isDomain {
		ip = FindDomainByIP(ip)
	}

	commonParams := FindCommonParamsByIp(ip)

	if reflect.DeepEqual(commonParams.User, xui.User{}) {
		commonParams.User = &xui.User{
			Password: "csh031027",
			UserName: "csh0101",
		}
	}

	var inbound *xui.Inbound

	remark = fmt.Sprintf("%s-%s", zone, remark)
	if !isTls {
		inbound = xui.NewVmessInbound(remark)
	} else {
		inbound = xui.NewVmessTLSInbound(remark)
	}

	_, err = xui.AddInbound(commonParams, inbound)
	if err != nil {
		return err
	}

	_, err = xui.AddOutbound(commonParams, &xui.Outbound{
		Protocol:    "freedom",
		SendThrough: ip,
		Tag:         "outbound-0",
	})

	if err != nil {
		return err
	}

	_, err = xui.AddRouterRule(commonParams, &xui.RouterRule{
		Type:        "field",
		InboundTag:  []string{"inbound-0"},
		OutboundTag: "outbound-0",
	})
	if err != nil {
		return err
	}

	return nil

}

func FindDomainByIP(ip string) string {
	return ""
}

func FindCommonParamsByIp(ip string) *xui.CommonParams {
	return &xui.CommonParams{
		Url: ip,
	}
}

func SeletIPByOrderID(db *gorm.DB, orderID string) ([]*models.Host, error) {
	var hostDeviceAllocates []models.HostDeviceAllocate
	db.Where("device_id = ?", orderID).Find(&hostDeviceAllocates)
	var hosts []*models.Host
	for _, allocate := range hostDeviceAllocates {
		host := &models.Host{}
		err := db.Where("id = ?", allocate.HostID).First(&host).Error
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

func AllocateForZone(db *gorm.DB, req *view.AllocateRequest, c echo.Context) error {

	hostDeviceAllocate := make([]*models.HostDeviceAllocate, 0)
	for zone, quota := range req.Lines {
		unallocatedHosts, err := getUnallocatedHosts(db, zone)
		if err != nil {
			return err
		}
		if len(unallocatedHosts) <= quota+1000 {
			return errors.New("IP not enough")
		}
		hosts := selectRandomHosts(unallocatedHosts, quota)
		for _, host := range hosts {
			hostDeviceAllocate = append(hostDeviceAllocate, &models.HostDeviceAllocate{
				HostID:  host.ID,
				OrderID: req.OrderID,
			})
		}
	}
	err := db.Create(hostDeviceAllocate).Error
	if err != nil {
		return err
	}
	return err
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
