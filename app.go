package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"strings"

	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/Hunter-club/cloudman/view"
	"github.com/Hunter-club/cloudman/xui"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	probing "github.com/prometheus-community/pro-bing"
	"gorm.io/gorm"
)

var SubURLPrefix string = "http://localhost:9999"

func GetHealthTransferProxy(transferProxy map[string]string) (string, error) {

	proxies := make([]string, 0)
	for proxy, port := range transferProxy {
		isHealth, err := PingHealth(proxy)
		if err != nil {
			continue
		}

		if isHealth {
			proxies = append(proxies, fmt.Sprintf("%s:%s", proxy, port))
		}
	}

	if len(proxies) == 0 {
		return "", errors.New("No Health Proxy")
	}

	index := rand.Intn(len(proxies))

	return proxies[index], nil
}

func PingHealth(addr string) (bool, error) {
	health := false
	p, err := probing.NewPinger(addr)
	if err != nil {
		return false, err
	}

	p.OnFinish = func(s *probing.Statistics) {
		if s.PacketLoss <= 0.3 {
			health = true
		}
	}
	p.Count = 5

	err = p.Run()
	if err != nil {
		return false, err
	}

	return health, nil
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := database.Db

	e.POST("/api/v1/sub/:order_id", func(c echo.Context) error {
		var err error
		req := &view.SubRequest{}
		err = c.Bind(req)
		if err != nil {
			return DealResp(c, "bind error", err, nil)
		}
		res := make([]string, 0)

		tx := db.Begin()
		defer func() {
			if err == nil {
				tx.Commit()
			} else {
				tx.Rollback()
			}
		}()

		orderID := c.Param("order_id")

		saveTransfer := make([]*models.Transfer, 0)
		hostTransfers := make([]*models.HostTransfer, 0)
		for _, entry := range req.Entries {
			commonParams := FindCommonParamsByIp(entry.IP)
			if reflect.DeepEqual(commonParams.User, xui.User{}) {
				commonParams.User = &xui.User{
					Password: "csh031027",
					UserName: "csh0101",
				}
			}
			inbound, err := xui.GetInbound(commonParams)
			if err != nil {
				DealResp(c, "error", err, nil)
			}
			vmess, err := GetSubJson(commonParams, xui.GetInboundSubId(inbound), entry.Transfer.Addr, entry.Transfer.Port)
			if err != nil {
				DealResp(c, "subjson error", err, nil)
			}
			res = append(res, vmess)

			host := &models.Host{}

			db.Model(&models.Host{}).Where(&models.Host{
				PrimaryIP: entry.IP,
			}).Select("host.*").Find(&host)

			transferID := uuid.NewString()
			saveTransfer = append(saveTransfer, &models.Transfer{
				Port:       entry.Transfer.Port,
				Addr:       entry.Transfer.Addr,
				TransferID: transferID,
			})

			hostTransfer := &models.HostTransfer{
				TransferID: transferID,
				HostID:     host.HostID,
			}
			hostTransfers = append(hostTransfers, hostTransfer)

		}

		err = db.Create(saveTransfer).Error
		if err != nil {
			return err
		}

		err = db.Create(hostTransfers).Error
		if err != nil {
			return err
		}

		vmessListSub := strings.Join(res, "\n")

		newSubID := uuid.NewString()[:16]
		orderSub := &models.OrderSub{
			OrderID: orderID,
			SubID:   newSubID,
			Vmess:   vmessListSub,
		}

		err = db.Create(orderSub).Error
		if err != nil {
			DealResp(c, "create orderSub", err, nil)
		}
		tempResp := struct {
			SubUrl string `json:"sub_url"`
		}{
			SubUrl: fmt.Sprintf("%s:%s:%s", SubURLPrefix, "/sub", newSubID),
		}
		return DealResp(c, "success", nil, tempResp)
	})

	e.POST("/api/v1/xray", func(c echo.Context) error {
		// 获取订阅的连接
		// 为设备创建固定的规则
		// 入站、出站、路由规则
		// 生成订阅
		var err error

		req := &view.XrayRequest{}

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

		hosts, err := SelectIPByOrderID(db, req.OrderID)
		if err != nil {
			return err
		}

		result := make([]xui.Inbound, 0)
		for _, host := range hosts {
			inbound, err := ConfigXraySingle(host, false, false)
			if err != nil {
				DealResp(c, "bind error", err, nil)
			}
			result = append(result, *inbound)
		}
		c.JSON(http.StatusOK, result)
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

func ConfigXraySingle(host *models.Host, isDomain, isTls bool) (*xui.Inbound, error) {
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

	remark := fmt.Sprintf("%s-%s", zone, host.PrimaryIP)
	if !isTls {
		inbound = xui.NewVmessInbound(remark)
	} else {
		inbound = xui.NewVmessTLSInbound(remark)
	}

	_, err = xui.AddInbound(commonParams, inbound)
	if err != nil {
		return nil, err
	}

	_, err = xui.AddOutbound(commonParams, &xui.Outbound{
		Protocol:    "freedom",
		SendThrough: ip,
		Tag:         "outbound-0",
	})
	if err != nil {
		return nil, err
	}

	_, err = xui.AddRouterRule(commonParams, &xui.RouterRule{
		Type:        "field",
		InboundTag:  []string{"inbound-0"},
		OutboundTag: "outbound-0",
	})
	if err != nil {
		return nil, err
	}

	return inbound, nil
}

func FindDomainByIP(ip string) string {
	return ""
}

func FindCommonParamsByIp(ip string) *xui.CommonParams {
	return &xui.CommonParams{
		Url: ip,
	}
}

func SelectIPByOrderID(db *gorm.DB, orderID string) ([]*models.Host, error) {
	var hostDeviceAllocates []models.HostDeviceAllocate
	db.Where(&models.HostDeviceAllocate{
		OrderID:    orderID,
		IsAllocate: false,
	}).Find(&hostDeviceAllocates)
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

func getAllocatedHosts(db *gorm.DB, orderID string) ([]models.Host, error) {
	var allocatedHosts []models.Host
	if err := db.Model(&models.Host{}).
		Joins("LEFT JOIN host_device_allocate ON host.id = host_device_allocate.host_id").
		Where("host_device_allocate.order_id = ?", orderID).
		Select("host.*").
		Find(&allocatedHosts).Error; err != nil {
		return nil, err
	}
	return allocatedHosts, nil
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

func DealResp(c echo.Context, msg string, err error, data interface{}) error {
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusOK, JsonObj{
			Msg:     msg,
			Success: false,
			Data:    data,
		})
	} else {
		return c.JSON(http.StatusOK, JsonObj{
			Msg:     msg,
			Success: true,
			Data:    data,
		})
	}
}

func GetSubJson(commonParams *xui.CommonParams, subId string, Addr string, Port int) (string, error) {
	subJson, err := xui.GetSubJson(commonParams, subId)
	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(subJson[8:len(subJson)]) // base64转换为json
	if err != nil {
		return "", err
	}

	var vmess view.Vmess
	err = json.Unmarshal([]byte(decoded), &vmess)
	if err != nil {
		fmt.Println("vmess convert error", err.Error())
	}
	vmess.Add = Addr
	vmess.Port = Port

	rawVmess, err := json.Marshal(vmess)
	if err != nil {
		return "", err
	}
	encode := base64.StdEncoding.EncodeToString(rawVmess) // byte转base64
	encode = "vmess://" + "" + encode                     // 加上vmess头

	return encode, nil
}
