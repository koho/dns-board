package controllers

import (
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/gin-gonic/gin"
	"github.com/koho/dns-board/db"
	"github.com/koho/dns-board/models"
	"net/http"
	"strings"
	"time"
)

type ClientStat struct {
	Requests   int                `json:"req"`
	Size       float64            `json:"-"`
	SizeInt    int                `json:"size"`
	Cache      int                `json:"cache"`
	RCode      map[int]int        `json:"rcode"`
	Forwarders map[string]float64 `json:"-"`
	FwdCounter map[string]int     `json:"-"`
	FwdInt     map[string]int     `json:"fwd"`
}

type Domain struct {
	IPv4  string `json:"ipv4"`
	IPv6  string `json:"ipv6"`
	Count int    `json:"count"`
}

func GetStat(c *gin.Context) {
	var msg []models.Message
	if err := models.AddTimeClause(db.GetDB(), c.DefaultQuery("hour", "3")).Find(&msg).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	cs := make(map[string]*ClientStat)
	ds := make(map[string]*Domain)
	ipc := make(map[string]int)
	qtc := make(map[string]int)
	dateTime := time.DateTime[:16]
	for _, m := range msg {
		minute := m.Time.Format(dateTime)
		s, ok := cs[minute]
		if !ok {
			s = &ClientStat{RCode: make(map[int]int), Forwarders: make(map[string]float64), FwdCounter: make(map[string]int), FwdInt: make(map[string]int)}
			cs[minute] = s
		}
		if m.Type == dnstap.Message_CLIENT_RESPONSE {
			s.Requests++
			if m.Duration == 0 {
				s.Cache++
			}
			s.Size -= (s.Size - float64(m.Size)) / float64(s.Requests)
			s.SizeInt = int(s.Size)
			s.RCode[m.RCode]++
			if m.Answer != "" {
				if m.QType == "A" || m.QType == "AAAA" {
					d, ok := ds[m.Domain]
					if !ok {
						d = &Domain{}
						ds[m.Domain] = d
					}
					ips := strings.Split(m.Answer, ",")
					if len(ips) > 0 {
						d.Count++
						switch m.QType {
						case "A":
							d.IPv4 = ips[0]
						case "AAAA":
							d.IPv6 = ips[0]
						}
					}
				}
				qtc[m.QType]++
				ipc[m.IP]++
			}
		} else if m.Type == dnstap.Message_FORWARDER_RESPONSE {
			s.FwdCounter[m.IP]++
			s.Forwarders[m.IP] -= (s.Forwarders[m.IP] - float64(m.Duration)) / float64(s.FwdCounter[m.IP])
			s.FwdInt[m.IP] = int(s.Forwarders[m.IP])
		}
	}
	c.JSON(http.StatusOK, gin.H{"cs": cs, "ds": ds, "ipc": ipc, "qtc": qtc})
}
