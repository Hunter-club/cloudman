package kits

import (
	"errors"
	"fmt"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func GetHealthTransferProxy(transferProxy map[string]string) (string, error) {

	proxies := make([]string, 0)
	for proxy, port := range transferProxy {
		isHealth, err := PingHealth(proxy)
		if err != nil {
			continue
		}

		if isHealth {
			proxies = append(proxies, fmt.Sprintf("%s:%s", proxy, port))
		}
	}

	if len(proxies) == 0 {
		return "", errors.New("no health proxy")
	}
	rander := GetRander()

	index := rander.Intn(len(proxies))

	return proxies[index], nil
}

func PingHealth(addr string) (bool, error) {
	health := false
	p, err := probing.NewPinger(addr)
	if err != nil {
		return false, err
	}

	p.Timeout = time.Second * 5

	p.OnFinish = func(s *probing.Statistics) {
		if s.PacketLoss <= 0.3 {
			health = true
		}
	}
	p.Count = 1

	err = p.Run()
	if err != nil {
		return false, err
	}

	return health, nil
}
