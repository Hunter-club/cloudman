package handler

import (
	"fmt"
	"os"
	"reflect"

	"github.com/Hunter-club/cloudman/config"
	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/Hunter-club/cloudman/view"
	"github.com/Hunter-club/cloudman/xui"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func XUIConfigure(c echo.Context) (interface{}, error) {
	var err error

	req := &view.XrayRequest{}

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

	hosts, err := SelectIPByOrderID(db, req.OrderID)
	if err != nil {
		return nil, err
	}

	result := make([]xui.Inbound, 0)
	for _, host := range hosts {
		inbound, err := ConfigXraySingle(host, false, false)
		if err != nil {
			return nil, err
		}
		result = append(result, *inbound)

		//todo 对于分配成功的XUI面板，要进行设置
		err = db.Model(&models.HostOrderAllocate{}).Where(&models.HostOrderAllocate{
			OrderID: req.OrderID,
			HostID:  host.HostID,
		}).Update("is_allocate", true).Error
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func SelectIPByOrderID(db *gorm.DB, orderID string) ([]*models.Host, error) {
	var hostDeviceAllocates []models.HostOrderAllocate
	db.Where(&models.HostOrderAllocate{
		OrderID:    orderID,
		IsAllocate: false,
	}).Find(&hostDeviceAllocates)
	var hosts []*models.Host
	// 找到机器
	for _, allocate := range hostDeviceAllocates {
		host := &models.Host{}
		err := db.Where("host_id = ?", allocate.HostID).First(&host).Error
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

func ConfigXraySingle(host *models.Host, isDomain, isTls bool) (*xui.Inbound, error) {
	var err error
	zone := host.Zone
	var url string
	if host.Domain != "" {
		url = config.Protocol + host.Domain
	} else {
		if os.Getenv("DEBUG") != "" {
			url = config.Protocol + "localhost" + ":" + config.Port
		} else {
			url = config.Protocol + host.PrimaryIP + ":" + config.Port
		}
	}

	commonParams := &xui.CommonParams{
		Url:  url,
		User: &xui.User{},
	}
	if reflect.DeepEqual(*commonParams.User, xui.User{}) {
		commonParams.User = &xui.User{
			Password: "csh031027",
			UserName: "csh0101",
		}
	}

	var inbound *xui.Inbound

	remark := fmt.Sprintf("%s-%s", zone, host.PrimaryIP)
	inbound = xui.NewVmessInbound(remark, isTls)

	_, err = xui.AddInbound(commonParams, inbound)
	if err != nil {
		return nil, err
	}

	outbound, err := xui.AddOutbound(commonParams, &xui.Outbound{
		Protocol:    "freedom",
		SendThrough: host.PrimaryIP,
		Tag:         fmt.Sprintf("outbound-%d", inbound.Port),
	})
	if err != nil {
		return nil, err
	}

	_, err = xui.AddRouterRule(commonParams, &xui.RouterRule{
		Type:        "field",
		InboundTag:  []string{inbound.Tag},
		OutboundTag: outbound.Tag,
	})
	if err != nil {
		return nil, err
	}

	return inbound, nil
}
