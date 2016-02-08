package common

import (
	"math/rand"
	"time"
)

// RandomDelay randomly sleeps for up to the configured number of milliseconds
func RandomDelay() {
	config, err := NewConfig()
	if err != nil {
		return
	}

	factor := config.DelayFactorMs
	if factor > 0 {
		time.Sleep(time.Duration(rand.Intn(factor)) * time.Millisecond)
	}
}
