package helpers

import (
	"math/rand"
	"time"
)

var (
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func GlobalRand() *rand.Rand {
	return globalRand
}
