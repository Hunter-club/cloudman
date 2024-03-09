package xui

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Hunter-club/cloudman/config"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
)

type CommonParams struct {
	Url  string
	User *User
}

func NewCommonParams(url string, user *User) *CommonParams {
	return &CommonParams{
		Url:  url,
		User: user,
	}
}

func (params *CommonParams) GetURL() string {
	return params.Url
}

func (params *CommonParams) GetUser() string {
	return params.Url
}

func Login(commonParams *CommonParams) ([]*http.Cookie, error) {
	req.DevMode()

	resp, err := req.C().NewRequest().SetBody(commonParams.User).Post(commonParams.Url + "/login")
	if err != nil {
		return nil, err
	}

	return resp.Cookies(), nil
}

func AddInbound(commonParams *CommonParams, inbound *Inbound) (bool, error) {
	req.DevMode()
	cookies, err := globalCookieTool.GetCookies(commonParams)
	if err != nil {
		return false, err
	}

	_, err = req.C().NewRequest().SetBody(inbound).SetCookies(cookies...).Post(commonParams.Url + "/xui/API/inbounds/add")
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetInbound(commonParams *CommonParams) (*Inbound, error) {
	req.DevMode()
	cookies, err := globalCookieTool.GetCookies(commonParams)
	if err != nil {
		return nil, err
	}
	resp, err := req.C().NewRequest().SetCookies(cookies...).Get(commonParams.Url + "/xui/API/inbounds")
	if err != nil {
		return nil, err
	}

	inbound := make([]*Inbound, 0)

	err = json.Unmarshal([]byte(gjson.GetBytes(resp.Bytes(), "obj").String()), &inbound)
	if err != nil {
		return nil, err
	}
	return inbound[0], nil
}

func AddOutbound(commonParams *CommonParams, outbound *Outbound) (*Outbound, error) {
	req.DevMode()

	cookies, err := globalCookieTool.GetCookies(commonParams)
	if err != nil {
		return outbound, err
	}

	_, err = req.C().NewRequest().SetBody(outbound).SetCookies(cookies...).Post(commonParams.Url + "/xui/API/outbounds/add")
	if err != nil {
		return outbound, err
	}
	return outbound, nil
}

func AddRouterRule(commonParams *CommonParams, routerRule *RouterRule) (bool, error) {
	req.DevMode()

	cookies, err := globalCookieTool.GetCookies(commonParams)
	if err != nil {
		return false, err
	}

	_, err = req.C().NewRequest().SetBody(routerRule).SetCookies(cookies...).Post(commonParams.Url + "/xui/API/routers/add")
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteRouterRule(commonParams *CommonParams, routerRule *RouterRule) (bool, error) {
	req.DevMode()
	cookies, err := globalCookieTool.GetCookies(commonParams)
	if err != nil {
		return false, err
	}
	_, err = req.C().NewRequest().SetBody(routerRule).SetCookies(cookies...).Post(commonParams.Url + "/xui/API/routers/delete")
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetSubJson(commonParams *CommonParams, subID string) (string, error) {
	req.DevMode()
	cookies, err := globalCookieTool.GetCookies(commonParams)
	if err != nil {
		return "", err
	}

	subURL := commonParams.Url[:strings.LastIndex(commonParams.Url, ":")] + ":" + config.SubPort
	resp, err := req.C().NewRequest().SetCookies(cookies...).SetPathParam("subID", subID).Get(subURL + "/sub/{subID}")
	if err != nil {
		return "", err
	}
	return string(resp.Bytes()), nil
}
