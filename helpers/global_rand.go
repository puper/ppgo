package helpers

import (
	"math/rand"
	"time"
)

func GlobalRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
