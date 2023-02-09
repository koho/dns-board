package server

import (
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/koho/dns-board/db"
	"github.com/koho/dns-board/models"
	"sync"
	"time"
)

var (
	identity = "DNS"
	version  = "0.0.0"
	startup  = time.Now().UnixMilli()
	mapURL   = ""
)

func init() {
	db.OnStartup(func(option db.Option) error {
		if id, err := models.GetMeta("identity"); err == nil && id != "" {
			identity = id
		}
		if ver, err := models.GetMeta("version"); err == nil && ver != "" {
			version = ver
		}
		return nil
	})
}

var meta sync.Once

func setMeta(dt *dnstap.Dnstap) {
	meta.Do(func() {
		identity = string(dt.Identity)
		version = string(dt.Version)
		models.SetMeta("identity", identity)
		models.SetMeta("version", version)
	})
}

type Option struct {
	Tap    string `yaml:"tap"`
	Listen string `yaml:"listen"`
	Map    string `yaml:"map"`
}

func Run(opt Option) {
	mapURL = opt.Map
	go RunDnstap(opt.Tap)
	RunWeb(opt.Listen)
}
