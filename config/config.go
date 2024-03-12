package config

import "os"

var SubURLPrefix string = "http://localhost:9999"
var Port = "54321"
var SubPort = "2096"
var Protocol string

func init() {
	Protocol = os.Getenv("PROTOCOL")
}
