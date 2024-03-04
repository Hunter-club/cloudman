package xui

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/tidwall/sjson"
)

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Inbound struct {
	Id          int             `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	UserId      int             `json:"-"`
	Up          int64           `json:"up" form:"up"`
	Down        int64           `json:"down" form:"down"`
	Total       int64           `json:"total" form:"total"`
	Remark      string          `json:"remark" form:"remark"`
	Enable      bool            `json:"enable" form:"enable"`
	ExpiryTime  int64           `json:"expiryTime" form:"expiryTime"`
	ClientStats []ClientTraffic `gorm:"foreignKey:InboundId;references:Id" json:"clientStats" form:"clientStats"`

	// config part
	Listen         string   `json:"listen" form:"listen"`
	Port           int      `json:"port" form:"port"`
	Protocol       Protocol `json:"protocol" form:"protocol"`
	Settings       string   `json:"settings" form:"settings"`
	StreamSettings string   `json:"streamSettings" form:"streamSettings"`
	Tag            string   `json:"tag" form:"tag" gorm:"unique"`
	Sniffing       string   `json:"sniffing" form:"sniffing"`
}

func NewVmessTLSInbound(Remark string) *Inbound {
	inboud := &Inbound{
		Remark:   Remark,
		Protocol: "vmess",
		Port:     rand.Intn(65535),
		Enable:   true,
	}

	inboud.Settings = GetInboundClient()
	inboud.StreamSettings = `{\n  \"network\": \"tcp\",\n  \"security\": \"tls\",\n  \"externalProxy\": [],\n  \"tlsSettings\": {\n    \"serverName\": \"\",\n    \"minVersion\": \"1.2\",\n    \"maxVersion\": \"1.3\",\n    \"cipherSuites\": \"\",\n    \"rejectUnknownSni\": false,\n    \"certificates\": [\n      {\n        \"certificateFile\": \"/root/pem.pem\",\n        \"keyFile\": \"/root/key.key\",\n        \"ocspStapling\": 3600\n      }\n    ],\n    \"alpn\": [\n      \"h2\",\n      \"http/1.1\"\n    ],\n    \"settings\": {\n      \"allowInsecure\": false,\n      \"fingerprint\": \"\"\n    }\n  },\n  \"tcpSettings\": {\n    \"acceptProxyProtocol\": false,\n    \"header\": {\n      \"type\": \"none\"\n    }\n  }\n}`
	inboud.Sniffing = `{\n  \"enabled\": true,\n  \"destOverride\": [\n    \"http\",\n    \"tls\",\n    \"quic\",\n    \"fakedns\"\n  ]\n}`

	return inboud
}

func NewVmessInbound(Remark string) *Inbound {
	inboud := &Inbound{
		Remark:   Remark,
		Protocol: "vmess",
		Enable:   true,
		Port:     rand.Intn(65535),
	}

	inboud.Settings = GetInboundClient()
	inboud.StreamSettings = "{\n  \"network\": \"tcp\",\n  \"security\": \"tls\",\n  \"externalProxy\": [],\n  \"tlsSettings\": {\n    \"serverName\": \"\",\n    \"minVersion\": \"1.2\",\n    \"maxVersion\": \"1.3\",\n    \"cipherSuites\": \"\",\n    \"rejectUnknownSni\": false,\n    \"certificates\": [\n      {\n        \"certificateFile\": \"/root/pem.pem\",\n        \"keyFile\": \"/root/key.key\",\n        \"ocspStapling\": 3600\n      }\n    ],\n    \"alpn\": [\n      \"h2\",\n      \"http/1.1\"\n    ],\n    \"settings\": {\n      \"allowInsecure\": false,\n      \"fingerprint\": \"\"\n    }\n  },\n  \"tcpSettings\": {\n    \"acceptProxyProtocol\": false,\n    \"header\": {\n      \"type\": \"none\"\n    }\n  }\n}"
	inboud.Sniffing = "{\n  \"enabled\": true,\n  \"destOverride\": [\n    \"http\",\n    \"tls\",\n    \"quic\",\n    \"fakedns\"\n  ]\n}"

	return inboud
}

func GetInboundClient() string {
	id := uuid.NewString()
	email := uuid.NewString()[:9]
	template := "{\n  \"clients\": [\n    {\n      \"id\": \"3aea0d2f-0fdd-424a-96e2-6d329a82c5a8\",\n      \"email\": \"9hybwkqp4\",\n      \"totalGB\": 0,\n      \"expiryTime\": 0,\n      \"enable\": true,\n      \"tgId\": \"\",\n      \"subId\": \"rrs3r602pg35j5om\",\n      \"reset\": 0\n    }\n  ]\n}"
	res, _ := sjson.Set(template, "clients.0.id", id)
	res, _ = sjson.Set(res, "clients.0.email", email)
	return res
}

type StreamSettings struct {
}

type ClientTraffic struct {
	Id         int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	InboundId  int    `json:"inboundId" form:"inboundId"`
	Enable     bool   `json:"enable" form:"enable"`
	Email      string `json:"email" form:"email" gorm:"unique"`
	Up         int64  `json:"up" form:"up"`
	Down       int64  `json:"down" form:"down"`
	ExpiryTime int64  `json:"expiryTime" form:"expiryTime"`
	Total      int64  `json:"total" form:"total"`
	Reset      int    `json:"reset" form:"reset" gorm:"default:0"`
}

type Protocol string

const (
	VMess       Protocol = "vmess"
	VLESS       Protocol = "vless"
	Dokodemo    Protocol = "Dokodemo-door"
	Http        Protocol = "http"
	Trojan      Protocol = "trojan"
	Shadowsocks Protocol = "shadowsocks"
)

type Outbound struct {
	Protocol    string `json:"protocol"`
	Settings    string `json:"settings"`
	Tag         string `json:"tag"`
	SendThrough string `json:"sendthrough"`
}

type RouterRule struct {
	Type        string   `json:"type"`
	InboundTag  []string `json:"inboundTag"`
	OutboundTag string   `json:"outboundTag"`
}
