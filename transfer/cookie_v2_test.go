package transfer_test

import (
	"testing"

	"github.com/Hunter-club/cloudman/transfer"
	"github.com/stretchr/testify/assert"
)

func TestTransferLogin(t *testing.T) {

	cookies, err := transfer.Login()

	assert.Nil(t, err)

	t.Log(cookies)
}
