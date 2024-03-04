package xui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCookiesUtils_Concurrency(t *testing.T) {
	url := "http://localhost:54321"
	cookiesUtils := NewCookiesUtils()

	// 模拟获取新Cookies的函
	// 验证获取到的Cookies是否正确
	cookies, err := cookiesUtils.GetCookies(&CommonParams{
		Url: url,
		User: &User{
			UserName: "csh0101",
			Password: "csh031027",
		},
	})
	assert.Nil(t, err)
	t.Log(cookies)
}
