package view

type Vmess struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Id   string `json:"id"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Tls  string `json:"tls"`
	Sni  string `json:"sni"`
	Alpn string `json:"alpn"`
	Port int    `json:"port"`
}
