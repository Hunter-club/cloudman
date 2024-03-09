package xui

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestXuiAddInbound(t *testing.T) {

	_, err := AddInbound(&CommonParams{
		Url: "http://localhost:54321",
		User: &User{
			UserName: "csh0101",
			Password: "csh031027",
		},
	}, NewVmessInbound("test1", false))

	assert.Nil(t, err)
}

func TestGetSubID(t *testing.T) {

	val := "{\n  \"clients\": [\n    {\n      \"id\": \"3d9381ed-6346-469e-b1bb-be26124b55c8\",\n      \"email\": \"a44391da-\",\n      \"totalGB\": 0,\n      \"expiryTime\": 0,\n      \"enable\": true,\n      \"tgId\": \"\",\n      \"subId\": \"ec6df146-fd17-49\",\n      \"reset\": 0\n    }\n  ]\n}"

	fmt.Println(gjson.Get(val, "clients.0.subId").String())
	fmt.Println(GetInboundSubId(&Inbound{
		Settings: val,
	}))
}
