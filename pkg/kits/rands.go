package kits

import (
	"math/rand"
	"time"
)

func GetRander() *rand.Rand {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rng
}
