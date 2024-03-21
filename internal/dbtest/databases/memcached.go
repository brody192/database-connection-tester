package databases

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func Memcache(server string) (time.Duration, error) {
	mc := memcache.New(server)

	sT := time.Now()

	if err := mc.Ping(); err != nil {
		return time.Since(sT), fmt.Errorf("mc.Ping error: %w", err)
	}

	return time.Since(sT), nil
}
