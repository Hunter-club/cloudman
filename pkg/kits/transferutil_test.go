package kits_test

import (
	"testing"

	"github.com/Hunter-club/cloudman/pkg/kits"
	"github.com/stretchr/testify/assert"
)

func TestPingHealth(t *testing.T) {
	flag, err := kits.PingHealth("14.128.37.18")

	assert.Nil(t, err)

	t.Log(flag)
}
