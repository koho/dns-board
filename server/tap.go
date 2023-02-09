package server

import (
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/koho/dns-board/db"
	"github.com/koho/dns-board/output"
	"log"
	"net"
	"os"
	"strings"
)

const (
	tcpPrefix  = "tcp://"
	unixPrefix = "unix://"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func RunDnstap(u string) {
	var i *dnstap.FrameStreamSockInput
	var err error
	if strings.HasPrefix(u, tcpPrefix) {
		addr := strings.TrimPrefix(u, tcpPrefix)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Fatalf("dnstap: failed to listen on %s: %v\n", addr, err)
		}
		i = dnstap.NewFrameStreamSockInput(l)
	} else if strings.HasPrefix(u, unixPrefix) {
		path := strings.TrimPrefix(u, unixPrefix)
		if i, err = dnstap.NewFrameStreamSockInputFromPath(path); err != nil {
			logger.Fatalf("dnstap: failed to open input socket %s: %v\n", path, err)
		}
	} else {
		logger.Fatalf("unsupported protocol: %v\n", u)
	}
	i.SetTimeout(0)
	i.SetLogger(logger)
	to := output.NewDBOutput(db.GetDB())
	to.SetCallback(setMeta)
	go i.ReadInto(to.GetOutputChannel())
	to.RunOutputLoop()
}
