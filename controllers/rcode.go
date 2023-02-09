package controllers

import (
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/gin-gonic/gin"
	"github.com/koho/dns-board/db"
	"github.com/koho/dns-board/models"
	"github.com/miekg/dns"
	"net/http"
)

type RCodeStat struct {
	T     string
	RCode int
	Value float64
}

func GetRCodeStat(c *gin.Context) {
	m := db.GetDB().Model(&models.Message{}).Select("strftime('%Y-%m-%d %H:%M', time, 'localtime') as t, r_code, count(*) as value").
		Where("type = ?", dnstap.Message_CLIENT_RESPONSE)
	m = models.AddTimeClause(m, c.DefaultQuery("hour", "3"))
	var rCodes []RCodeStat
	if err := m.Group("t, r_code").Find(&rCodes).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	codeMap := make(map[int]bool)
	for _, r := range rCodes {
		codeMap[r.RCode] = false
	}
	result := make([]Record, 0, len(rCodes))
	addAbsent := func(t string) {
		if t != "" {
			for code, added := range codeMap {
				if !added {
					result = append(result, Record{t, dns.RcodeToString[code], 0})
				}
				codeMap[code] = false
			}
		}
	}
	prevTime := ""
	for _, r := range rCodes {
		if r.T != prevTime {
			addAbsent(prevTime)
		}
		result = append(result, Record{r.T, dns.RcodeToString[r.RCode], r.Value})
		codeMap[r.RCode] = true
		prevTime = r.T
	}
	addAbsent(prevTime)
	c.JSON(http.StatusOK, result)
}
