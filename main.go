package main

import (
	"flag"
	"fmt"
	"github.com/koho/dnstap-web/db"
	"github.com/koho/dnstap-web/models"
	"github.com/koho/dnstap-web/server"
	"log"
)

const version = "1.0.0"

var (
	tapUrl      string
	webListen   string
	setPass     string
	showVersion bool
)

func init() {
	flag.StringVar(&tapUrl, "t", "tcp://:6000", "dnstap url (tcp or unix).")
	flag.StringVar(&webListen, "l", ":80", "web server listen address")
	flag.StringVar(&server.MapURL, "m", "", "maplibre style url")
	flag.StringVar(&setPass, "p", "", "set admin password")
	flag.BoolVar(&showVersion, "v", false, "show version info")
	flag.Parse()
}

func main() {
	if showVersion {
		fmt.Println(version)
		return
	}
	db.Init()
	if setPass != "" {
		if err := models.UpdateUserPassword("admin", setPass); err != nil {
			log.Fatal(err)
		}
		return
	}
	go server.RunDnstap(tapUrl)
	server.RunWeb(webListen)
}
