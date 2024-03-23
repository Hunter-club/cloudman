package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/Hunter-club/cloudman/config"
	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/Hunter-club/cloudman/pkg/kits"
	"github.com/Hunter-club/cloudman/view"
	"github.com/Hunter-club/cloudman/xui"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GenSub(c echo.Context) (interface{}, error) {
	fmt.Println("gen sub is processing")
	var err error
	req := &view.SubRequest{}
	err = c.Bind(req)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0)

	db := database.GetDB()

	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	orderID := req.OrderID

	saveTransfer := make([]*models.Transfer, 0)
	for _, entry := range req.Entries {
		host := &models.Host{}
		db.Model(&models.Host{}).Where(&models.Host{
			PrimaryIP: entry.TargetHost.Addr,
		}).Find(&host)
		//如果域名不对
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
		inbound, err := xui.GetInbound(commonParams)
		if err != nil {
			return nil, err
		}

		// 从运营增加的中转中挑出合适的中转

		transferMap := make(map[string]string)

		for _, v := range entry.Transfer {
			transferMap[v.Addr] = strconv.Itoa(v.Port)
			transferID := uuid.NewString()
			saveTransfer = append(saveTransfer, &models.Transfer{
				TransferID: transferID,
				Addr:       v.Addr,
				Port:       v.Port,
				HostID:     host.HostID,
				OrderID:    orderID,
				HostIP:     host.PrimaryIP,
			})

		}

		healthProxy, err := kits.GetHealthTransferProxy(transferMap)

		if err != nil {
			return nil, err
		}

		if len(healthProxy) == 0 {
			return nil, errors.New("no healthProxy")
		}

		selectedProxy := strings.Split(healthProxy, ":")
		selectPort, err := strconv.Atoi(selectedProxy[1])

		if err != nil {
			return nil, err
		}

		vmess, vmessID, err := GetSubJson(commonParams, xui.GetInboundSubId(inbound), selectedProxy[0], selectPort)
		if err != nil {
			return nil, err
		}
		res = append(res, vmess)

		// 把订阅链接保存到对应的中转 修改对应的订阅时，要找到可用的中转
		for _, transfer := range saveTransfer {
			transfer.VmessID = vmessID
		}

		err = db.Create(saveTransfer).Error
		if err != nil {
			return nil, err
		}
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
		return nil, err
	}
	tempResp := struct {
		SubUrl string `json:"sub_url"`
	}{
		SubUrl: fmt.Sprintf("%s/%s/%s", config.SubURLPrefix, "sub", newSubID),
	}
	return tempResp, nil
}

func GetSubJson(commonParams *xui.CommonParams, subId string, Addr string, Port int) (string, string, error) {
	subJson, err := xui.GetSubJson(commonParams, subId)
	if err != nil {
		return "", "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(subJson) // base64转换为json
	if err != nil {
		return "", "", err
	}

	decoded = decoded[8:]

	var vmess view.Vmess

	decoded, err = base64.StdEncoding.DecodeString(string(decoded)) // base64转换为json
	if err != nil {
		return "", "", err
	}
	err = json.Unmarshal([]byte(decoded), &vmess)
	if err != nil {
		fmt.Println("vmess convert error", err.Error())
	}
	vmess.Add = Addr
	vmess.Port = Port

	rawVmess, err := json.Marshal(vmess)
	if err != nil {
		return "", "", err
	}
	encode := base64.StdEncoding.EncodeToString(rawVmess) // byte转base64
	encode = "vmess://" + "" + encode                     // 加上vmess头

	return encode, vmess.Id, nil
}

// 定时订阅的
func CheckSubJsonV2() error {

	db := database.GetDB()

	orderSubs := make([]*models.OrderSub, 0)

	err := db.Model(&models.OrderSub{}).Find(&orderSubs).Error
	if err != nil {
		return err
	}

	for _, orderSub := range orderSubs {

		//拿到完整的vmess
		subJson := orderSub.Vmess
		newVmessResults := make([]string, 0)

		oldVmessSubJsons := strings.Split(string(subJson), "\n")

		for _, subJson := range oldVmessSubJsons {
			subJson = subJson[8:]
			subJson, err := base64.StdEncoding.DecodeString(string(subJson))
			if err != nil {
				return err
			}
			var vmess view.Vmess
			err = json.Unmarshal([]byte(subJson), &vmess)
			if err != nil {
				fmt.Println("vmess convert error", err.Error())
			}

			vmessTransfer := make([]models.Transfer, 0)

			db.Model(&models.Transfer{}).Where(&models.Transfer{
				VmessID: vmess.Id,
			}).Find(&vmessTransfer)

			transferMap := make(map[string]string)

			for _, v := range vmessTransfer {
				transferMap[v.Addr] = strconv.Itoa(v.Port)
			}

			healthProxy, err := kits.GetHealthTransferProxy(transferMap)
			if err != nil {
				return err
			}

			if len(healthProxy) == 0 {
				return errors.New("no healthProxy")
			}

			selectedProxy := strings.Split(healthProxy, ":")
			selectPort, err := strconv.Atoi(selectedProxy[1])

			if err != nil {
				return err
			}
			//从Tag中拿到单挑订阅的IP us-182.382.42.3
			vmess.Add = selectedProxy[0]
			vmess.Port = selectPort
			//设置好可用的订阅

			rawVmess, err := json.Marshal(vmess)
			if err != nil {
				return err
			}
			encode := base64.StdEncoding.EncodeToString(rawVmess)
			encode = "vmess://" + "" + encode
			newVmessResults = append(newVmessResults, encode)
		}

		// 更新最新的订阅链接
		newRawVmess := strings.Join(newVmessResults, "/n")
		err = db.Model(&models.OrderSub{}).Where(&models.OrderSub{
			OrderID: orderSub.OrderID,
		}).Update("vmess", newRawVmess).Error
		if err != nil {
			return err
		}
	}
	return nil
}
