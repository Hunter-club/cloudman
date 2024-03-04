package xui

import (
	"fmt"
	"testing"
)

func TestXuiClientSetting(t *testing.T) {

	client := GetInboundClient()

	fmt.Println(client)

}
