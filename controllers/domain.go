package controllers

import (
	"github.com/deckarep/golang-set/v2"
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/gin-gonic/gin"
	"github.com/koho/dnstap-web/models"
	"github.com/samber/lo"
	"strings"
)

type DomainStat struct {
	Name  string             `json:"name"`
	IPv4  mapset.Set[string] `json:"ipv4"`
	IPv6  mapset.Set[string] `json:"ipv6"`
	Count uint64             `json:"count"`
}

func GetDomainTable(c *gin.Context) {
	messages, err := models.GetMessagesByType(dnstap.Message_CLIENT_RESPONSE, c.DefaultQuery("hour", "3"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	table := make(map[string]*DomainStat)
	for _, msg := range messages {
		if msg.Answer == "" {
			continue
		}
		getStat := func() (*DomainStat, mapset.Set[string]) {
			stat, ok := table[msg.Domain]
			if !ok {
				stat = &DomainStat{IPv4: mapset.NewSet[string](), IPv6: mapset.NewSet[string]()}
				table[msg.Domain] = stat
			}
			newIP := mapset.NewSet(strings.Split(msg.Answer, ",")...)
			if newIP.Cardinality() > 0 {
				stat.Count++
			}
			return stat, newIP
		}
		switch msg.QType {
		case "A":
			stat, newIP := getStat()
			stat.IPv4 = stat.IPv4.Union(newIP)
		case "AAAA":
			stat, newIP := getStat()
			stat.IPv6 = stat.IPv6.Union(newIP)
		}
	}
	domainList := lo.MapToSlice(table, func(key string, value *DomainStat) DomainStat {
		return DomainStat{
			Name:  key,
			IPv4:  value.IPv4,
			IPv6:  value.IPv6,
			Count: value.Count,
		}
	})
	c.JSON(200, domainList)
}
