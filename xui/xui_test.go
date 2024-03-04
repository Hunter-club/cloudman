package xui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXuiAddInbound(t *testing.T) {

	_, err := AddInbound(&CommonParams{
		Url: "http://localhost:54321",
		User: &User{
			UserName: "csh0101",
			Password: "csh031027",
		},
	}, NewVmessInbound("test1"))

	assert.Nil(t, err)
}
