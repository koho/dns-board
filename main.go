package main

import (
	"flag"
	"fmt"
	"github.com/koho/dns-board/db"
	"github.com/koho/dns-board/models"
	"github.com/koho/dns-board/server"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const version = "1.0.0"

var config struct {
	DB     db.Option     `yaml:",inline"`
	Server server.Option `yaml:",inline"`
}

var (
	confPath    string
	setPass     string
	showVersion bool
)

func init() {
	flag.StringVar(&confPath, "c", "config.yml", "config file path")
	flag.StringVar(&config.Server.Tap, "t", "tcp://:6000", "dnstap url (tcp or unix).")
	flag.StringVar(&config.Server.Listen, "l", ":80", "web server listen address")
	flag.StringVar(&config.Server.Map, "m", "", "maplibre style url")
	flag.StringVar(&setPass, "p", "", "set admin password")
	flag.IntVar(&config.DB.Retention, "r", 7, "the max number of day for data to keep")
	flag.BoolVar(&showVersion, "v", false, "show version info")
	flag.Parse()
}

func main() {
	if showVersion {
		fmt.Println(version)
		return
	}
	if cb, err := os.ReadFile(confPath); err == nil {
		if err = yaml.Unmarshal(cb, &config); err != nil {
			log.Fatal(err)
		}
	}

	db.Init(config.DB)
	if setPass != "" {
		if err := models.UpdateUserPassword("admin", setPass); err != nil {
			log.Fatal(err)
		}
		return
	}
	server.Run(config.Server)
}
