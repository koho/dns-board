package server

import (
	dnstap "github.com/dnstap/golang-dnstap"
	"sync"
	"time"
)

var (
	Identity = "DNS"
	Version  = "0.0.0"
	Startup  = time.Now().UnixMilli()
	MapURL   = ""
)

var meta sync.Once

func setMeta(dt *dnstap.Dnstap) {
	meta.Do(func() {
		Identity = string(dt.Identity)
		Version = string(dt.Version)
	})
}
