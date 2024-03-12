package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/Hunter-club/cloudman/config"
	"github.com/Hunter-club/cloudman/database"
	"github.com/Hunter-club/cloudman/models"
	"github.com/Hunter-club/cloudman/view"
	"github.com/Hunter-club/cloudman/xui"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GenSub(c echo.Context) (interface{}, error) {

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
	hostTransfers := make([]*models.HostTransfer, 0)
	for _, entry := range req.Entries {

		host := &models.Host{}

		db.Model(&models.Host{}).Where(&models.Host{
			PrimaryIP: entry.TargetHost.Addr,
		}).Find(&host)

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
		vmess, err := GetSubJson(commonParams, xui.GetInboundSubId(inbound), entry.Transfer.Addr, entry.Transfer.Port)
		if err != nil {
			return nil, err
		}
		res = append(res, vmess)

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
		return nil, err
	}

	err = db.Create(hostTransfers).Error
	if err != nil {
		return nil, err
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

func GetSubJson(commonParams *xui.CommonParams, subId string, Addr string, Port int) (string, error) {
	subJson, err := xui.GetSubJson(commonParams, subId)
	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(subJson) // base64转换为json
	if err != nil {
		return "", err
	}

	decoded = decoded[8:]

	var vmess view.Vmess

	decoded, err = base64.StdEncoding.DecodeString(string(decoded)) // base64转换为json
	if err != nil {
		return "", err
	}
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
