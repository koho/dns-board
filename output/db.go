package output

import (
	"errors"
	"fmt"
	"github.com/dnstap/golang-dnstap"
	"github.com/koho/dnstap-web/models"
	"github.com/miekg/dns"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type DBOutput struct {
	db            *gorm.DB
	outputChannel chan []byte
	cb            func(*dnstap.Dnstap)
	wait          chan bool
}

func NewDBOutput(db *gorm.DB) *DBOutput {
	return &DBOutput{
		db:            db,
		outputChannel: make(chan []byte, 32),
		wait:          make(chan bool),
	}
}

func (d *DBOutput) GetOutputChannel() chan []byte {
	return d.outputChannel
}

func (d *DBOutput) RunOutputLoop() {
	dt := &dnstap.Dnstap{}
	for frame := range d.outputChannel {
		if err := proto.Unmarshal(frame, dt); err != nil {
			log.Printf("proto.Unmarshal() failed: %s", err)
			continue
		}
		if d.cb != nil {
			d.cb(dt)
		}
		msg, err := d.MessageFromDnstap(dt.Message)
		if err != nil {
			if !errors.Is(err, os.ErrInvalid) {
				log.Printf("decode message failed: %s", err)
			}
			continue
		}
		if err = d.db.Create(msg).Error; err != nil {
			log.Printf("create message failed: %s", err)
		}
	}
	close(d.wait)
}

func (d *DBOutput) MessageFromDnstap(m *dnstap.Message) (*models.Message, error) {
	var ip string
	var port uint32
	switch m.GetType() {
	case dnstap.Message_CLIENT_RESPONSE:
		ip = net.IP(m.QueryAddress).String()
		port = m.GetQueryPort()
	case dnstap.Message_FORWARDER_RESPONSE:
		ip = net.IP(m.ResponseAddress).String()
		port = m.GetResponsePort()
	default:
		return nil, os.ErrInvalid
	}
	if m.ResponseMessage == nil {
		return nil, fmt.Errorf("invalid response message")
	}
	msg := new(dns.Msg)
	if err := msg.Unpack(m.ResponseMessage); err != nil {
		return nil, err
	}
	if len(msg.Question) == 0 {
		return nil, fmt.Errorf("empty question list")
	}
	record := models.Message{
		Time:     parseTime(m.GetQueryTimeSec(), m.GetQueryTimeNsec()),
		Type:     m.GetType(),
		IP:       ip,
		Port:     port,
		Protocol: m.GetSocketProtocol().String(),
		Domain:   msg.Question[0].Name,
		QType:    dns.TypeToString[msg.Question[0].Qtype],
		Answer:   gatherAnswer(msg.Question[0].Qtype, msg.Answer),
		RCode:    msg.Rcode,
	}
	respTime := parseTime(m.GetResponseTimeSec(), m.GetResponseTimeNsec())
	if record.Time != nil && respTime != nil {
		record.Duration = respTime.Sub(*record.Time).Milliseconds()
	}
	return &record, nil
}

func (d *DBOutput) SetCallback(cb func(*dnstap.Dnstap)) {
	d.cb = cb
}

func parseTime(sec uint64, nsec uint32) *time.Time {
	if sec != 0 && nsec != 0 {
		t := time.Unix(int64(sec), int64(nsec)).Local()
		return &t
	}
	return nil
}

func gatherAnswer(qtype uint16, ans []dns.RR) string {
	answers := make([]string, 0)
	for _, an := range ans {
		rrType := an.Header().Rrtype
		if rrType != qtype && (qtype == dns.TypeA || qtype == dns.TypeAAAA) {
			continue
		}
		var answer string
		switch rrType {
		case dns.TypeA:
			answer = an.(*dns.A).A.String()
		case dns.TypeAAAA:
			answer = an.(*dns.AAAA).AAAA.String()
		case dns.TypeCNAME:
			answer = an.(*dns.CNAME).Target
		case dns.TypeNS:
			answer = an.(*dns.NS).Ns
		case dns.TypePTR:
			answer = an.(*dns.PTR).Ptr
		case dns.TypeTXT:
			answer = strings.Join(an.(*dns.TXT).Txt, ",")
		case dns.TypeMX:
			answer = an.(*dns.MX).Mx
		case dns.TypeSOA:
			answer = an.(*dns.SOA).String()
		case dns.TypeSRV:
			answer = an.(*dns.SRV).String()
		case dns.TypeHTTPS:
			answer = an.(*dns.HTTPS).String()
		default:
			answer = an.String()
		}
		answers = append(answers, answer)
	}
	return strings.Join(answers, ",")
}

func (d *DBOutput) Close() {
	close(d.outputChannel)
	<-d.wait
}
