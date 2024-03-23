package config

import "os"

var SubURLPrefix string = "http://localhost:9999"
var Port = "54321"
var SubPort = "2096"
var Protocol string

func init() {
	if os.Getenv("PROTOCOL") != "" {
		Protocol = os.Getenv("PROTOCOL")
	}
	if os.Getenv("SUB_URL_PREFIX") != "" {
		SubURLPrefix = os.Getenv("SUB_URL_PREFIX")
	}
	if os.Getenv("SUB_PORT") != "" {
		SubPort = os.Getenv("SUB_PORT")
	}
	if os.Getenv("PORT") != "" {
		Port = os.Getenv("PORT")
	}
}
