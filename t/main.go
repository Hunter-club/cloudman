package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type Protocol string

const (
	ProtocolsVmess Protocol = "vmess"
)

type StreamSettings struct {
	Network  string
	Security string
	TCP      TCPSettings
	KCP      KCPSettings
	WS       WSSettings
	HTTP     HTTPSettings
	QUIC     QUICSettings
	GRPC     GRPCSettings
	TLS      TLSSettings
}

type TCPSettings struct {
	Type    string
	Request TCPRequest
}

type TCPRequest struct {
	Path    []string
	Headers []TCPHeader
}

type TCPHeader struct {
	Name  string
	Value string
}

type KCPSettings struct {
	Type string
	Seed string
}

type WSSettings struct {
	Path    string
	Headers []TCPHeader // Assuming similar structure to TCP for simplification
}

type HTTPSettings struct {
	Path string
	Host []string
}

type QUICSettings struct {
	Type     string
	Security string
	Key      string
}

type GRPCSettings struct {
	ServiceName string
	MultiMode   bool
}

type TLSSettings struct {
	SNI      string
	Settings TLSInnerSettings
	ALPN     []string
}

type TLSInnerSettings struct {
	Fingerprint   string
	AllowInsecure bool
}

// Helper function to check if a string is empty
func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Helper function to find header and return its value
func findHeader(headers []TCPHeader, name string) string {
	for _, header := range headers {
		if strings.ToLower(header.Name) == strings.ToLower(name) {
			return header.Value
		}
	}
	return ""
}

func genVmessLink(address string, port int, forceTls string, remark string, clientId string, stream StreamSettings, protocol Protocol) string {
	if protocol != ProtocolsVmess {
		return ""
	}

	security := forceTls
	if forceTls == "same" {
		security = stream.Security
	}

	obj := map[string]interface{}{
		"v":    "2",
		"ps":   remark,
		"add":  address,
		"port": port,
		"id":   clientId,
		"net":  stream.Network,
		"type": "none",
		"tls":  security,
	}
	fmt.Println("stream.Network:", stream.Network)

	switch stream.Network {
	case "tcp":
		obj["type"] = stream.TCP.Type
		if stream.TCP.Type == "http" {
			obj["path"] = strings.Join(stream.TCP.Request.Path, ",")
			obj["host"] = findHeader(stream.TCP.Request.Headers, "host")
		}
	case "kcp":
		obj["type"] = stream.KCP.Type
		obj["path"] = stream.KCP.Seed
	case "ws":
		obj["path"] = stream.WS.Path
		obj["host"] = findHeader(stream.WS.Headers, "host")
	case "http":
		obj["net"] = "h2"
		obj["path"] = stream.HTTP.Path
		obj["host"] = strings.Join(stream.HTTP.Host, ",")
	case "quic":
		obj["type"] = stream.QUIC.Type
		obj["host"] = stream.QUIC.Security
		obj["path"] = stream.QUIC.Key
	case "grpc":
		obj["path"] = stream.GRPC.ServiceName
		if stream.GRPC.MultiMode {
			obj["type"] = "multi"
		}
	}

	if security == "tls" {
		if !isEmpty(stream.TLS.SNI) {
			obj["sni"] = stream.TLS.SNI
		}
		if !isEmpty(stream.TLS.Settings.Fingerprint) {
			obj["fp"] = stream.TLS.Settings.Fingerprint
		}
		if len(stream.TLS.ALPN) > 0 {
			obj["alpn"] = strings.Join(stream.TLS.ALPN, ",")
		}
		if stream.TLS.Settings.AllowInsecure {
			obj["allowInsecure"] = stream.TLS.Settings.AllowInsecure
		}
	}

	fmt.Println(obj)
	jsonBytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ""
	}
	return "vmess://" + base64.StdEncoding.EncodeToString(jsonBytes)
}

type Config struct {
	ID             int          `json:"id"`
	Up             int64        `json:"up"`
	Down           int64        `json:"down"`
	Total          int          `json:"total"`
	Remark         string       `json:"remark"`
	Enable         bool         `json:"enable"`
	ExpiryTime     int64        `json:"expiryTime"`
	ClientStats    []ClientStat `json:"clientStats"`
	Listen         string       `json:"listen"`
	Port           int          `json:"port"`
	Protocol       string       `json:"protocol"`
	Settings       string       `json:"settings"`
	StreamSettings string       `json:"streamSettings"` // This will be parsed separately
	Tag            string       `json:"tag"`
	Sniffing       string       `json:"sniffing"`
}

type ClientStat struct {
	ID         int    `json:"id"`
	InboundID  int    `json:"inboundId"`
	Enable     bool   `json:"enable"`
	Email      string `json:"email"`
	Up         int64  `json:"up"`
	Down       int64  `json:"down"`
	ExpiryTime int64  `json:"expiryTime"`
	Total      int    `json:"total"`
	Reset      int    `json:"reset"`
}

func main() {
	// 假设 jsonStr 是包含上述 JSON 结构的字符串
	jsonStr := `
	     {
            "id": 2,
            "up": 2270587144,
            "down": 84673996961,
            "total": 0,
            "remark": "154.21.194.96",
            "enable": true,
            "expiryTime": 0,
            "clientStats": [
                {
                    "id": 2,
                    "inboundId": 2,
                    "enable": true,
                    "email": "8imrrdeqq",
                    "up": 2018078537,
                    "down": 82915187271,
                    "expiryTime": 0,
                    "total": 0,
                    "reset": 0
                }
            ],
            "listen": "",
            "port": 59296,
            "protocol": "vmess",
            "settings": "{\n  \"clients\": [\n    {\n      \"id\": \"adedd81d-5150-4db5-d5de-6f2f9920f083\",\n      \"email\": \"8imrrdeqq\",\n      \"totalGB\": 0,\n      \"expiryTime\": 0,\n      \"enable\": true,\n      \"tgId\": \"\",\n      \"subId\": \"hpn8jak89p0cg0d8\",\n      \"reset\": 0\n    }\n  ]\n}",
            "streamSettings": "{\n  \"network\": \"tcp\",\n  \"security\": \"tls\",\n  \"externalProxy\": [],\n  \"tlsSettings\": {\n    \"serverName\": \"isp.fanayun.com\",\n    \"minVersion\": \"1.2\",\n    \"maxVersion\": \"1.3\",\n    \"cipherSuites\": \"\",\n    \"rejectUnknownSni\": false,\n    \"certificates\": [\n      {\n        \"certificateFile\": \"/root/ssl/pem.pem\",\n        \"keyFile\": \"/root/ssl/key.key\",\n        \"ocspStapling\": 3600\n      }\n    ],\n    \"alpn\": [\n      \"h2\",\n      \"http/1.1\"\n    ],\n    \"settings\": {\n      \"allowInsecure\": false,\n      \"fingerprint\": \"\"\n    }\n  },\n  \"tcpSettings\": {\n    \"acceptProxyProtocol\": false,\n    \"header\": {\n      \"type\": \"none\"\n    }\n  }\n}",
            "tag": "inbound-59296",
            "sniffing": "{\n  \"enabled\": true,\n  \"destOverride\": [\n    \"http\",\n    \"tls\",\n    \"quic\",\n    \"fakedns\"\n  ]\n}"
        }

	` // 将这里替换为实际的 JSON 字符串

	var config Config
	if err := json.Unmarshal([]byte(jsonStr), &config); err != nil {
		fmt.Println("Error parsing config JSON:", err)
		return
	}

	var streamSettings StreamSettings
	if err := json.Unmarshal([]byte(config.StreamSettings), &streamSettings); err != nil {
		fmt.Println("Error parsing streamSettings JSON:", err)
		return
	}

	// 假设其他参数已经正确设置
	// 调用 genVmessLink 函数生成链接
	remark := config.Remark
	clientId := "adedd81d-5150-4db5-d5de-6f2f9920f083" // 需要从 config.Settings 或其他地方获取
	address := "154.21.194.96"                         // 需要根据具体情况设置
	forceTls := "same"                                 // 需要根据具体情况设置

	link := genVmessLink(address, config.Port, forceTls, remark, clientId, streamSettings, Protocol(config.Protocol))
	fmt.Println(link)
}
