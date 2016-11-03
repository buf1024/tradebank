package util

import (
	"fmt"
	"math/rand"
	"time"
)

func SID() string {
	return fmt.Sprintf("SIM-%d-%d",
		time.Now().Unix(), rand.Int63())
}

func init() {
	rand.NewSource(time.Now().Nanosecond())
}
