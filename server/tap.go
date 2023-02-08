package server

import (
	"fmt"
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/koho/dnstap-web/db"
	"github.com/koho/dnstap-web/output"
	"log"
	"net"
	"net/url"
	"os"
)

var logger = log.New(os.Stderr, "", log.LstdFlags)

func RunDnstap(u string) {
	tapUrl, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}
	var i *dnstap.FrameStreamSockInput
	switch tapUrl.Scheme {
	case "tcp":
		l, err := net.Listen("tcp", tapUrl.Host)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dnstap: Failed to listen on %s: %v\n", tapUrl.Host, err)
			os.Exit(1)
		}
		i = dnstap.NewFrameStreamSockInput(l)
	case "unix":
		i, err = dnstap.NewFrameStreamSockInputFromPath(tapUrl.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dnstap: Failed to open input socket %s: %v\n", tapUrl.Path, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "dnstap: opened input socket %s\n", tapUrl.Path)
	}
	i.SetTimeout(0)
	i.SetLogger(logger)
	to := output.NewDBOutput(db.GetDB())
	to.SetCallback(setMeta)
	go i.ReadInto(to.GetOutputChannel())
	to.RunOutputLoop()
}
