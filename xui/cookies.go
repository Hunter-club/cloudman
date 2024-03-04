package xui

import (
	"net/http"
	"sync"
)

var globalCookieTool *CookiesUtils

func init() {
	globalCookieTool = NewCookiesUtils()
}

type CookiesUtils struct {
	hostCookies map[string][]*http.Cookie
	lock        sync.RWMutex
}

func NewCookiesUtils() *CookiesUtils {
	return &CookiesUtils{
		hostCookies: make(map[string][]*http.Cookie),
	}
}

func (u *CookiesUtils) GetCookies(commonParams *CommonParams) ([]*http.Cookie, error) {

	url := commonParams.Url

	u.lock.RLock()
	if cookies, ok := u.hostCookies[url]; ok {
		u.lock.RUnlock()
		return cookies, nil
	}
	u.lock.RUnlock()

	// 如果未找到Cookies，尝试刷新或获取新的Cookies
	u.lock.Lock()
	defer u.lock.Unlock()

	// 双重检查，以防在获取写锁的过程中其他线程已经更新了Cookies
	if cookies, ok := u.hostCookies[url]; ok {
		return cookies, nil
	}
	// 模拟获取新的Cookies过程
	newCookies, err := Login(commonParams)
	if err != nil {
		return nil, err
	}
	u.hostCookies[url] = newCookies
	return newCookies, nil
}
