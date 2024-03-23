package transfer

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var transferCookieTool *CookiesUtilsV2
var urlx = "https://pfgo.cmze.one/"

func init() {
	transferCookieTool = NewCookiesUtils()
}

type CookiesUtilsV2 struct {
	hostCookies map[string][]*http.Cookie
	lock        sync.RWMutex
}

func NewCookiesUtils() *CookiesUtilsV2 {
	return &CookiesUtilsV2{
		hostCookies: make(map[string][]*http.Cookie),
	}
}

func (u *CookiesUtilsV2) GetCookies() ([]*http.Cookie, error) {
	u.lock.RLock()
	if cookies, ok := u.hostCookies[urlx]; ok {
		u.lock.RUnlock()
		return cookies, nil
	}
	u.lock.RUnlock()

	// 如果未找到Cookies，尝试刷新或获取新的Cookies
	u.lock.Lock()
	defer u.lock.Unlock()

	// 双重检查，以防在获取写锁的过程中其他线程已经更新了Cookies
	if cookies, ok := u.hostCookies[urlx]; ok {
		return cookies, nil
	}
	// 模拟获取新的Cookies过程
	newCookies, err := Login()
	if err != nil {
		return nil, err
	}
	u.hostCookies[urlx] = newCookies
	return newCookies, nil
}

func Login() ([]*http.Cookie, error) {
	urlx := "https://pfgo.cmze.one/ajax/login"
	method := "POST"

	data := url.Values{
		"username": {"admin"},
		"password": {"admin"},
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, urlx, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	return resp.Cookies(), nil

}
