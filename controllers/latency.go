package controllers

import (
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/gin-gonic/gin"
	"github.com/koho/dns-board/db"
	"github.com/koho/dns-board/models"
	"net/http"
)

type Duration struct {
	T      string  `json:"time"`
	Server string  `json:"type"`
	Value  float64 `json:"value"`
}

func GetRequestDuration(c *gin.Context) {
	m := db.GetDB().Model(&models.Message{}).Select("strftime('%Y-%m-%d %H:%M', time, 'localtime') as t, ip as server, round(avg(duration), 0) as value").
		Where("type = ?", dnstap.Message_FORWARDER_RESPONSE)
	m = models.AddTimeClause(m, c.DefaultQuery("hour", "3"))
	var durations []Duration
	if err := m.Group("t, ip").Find(&durations).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, durations)
}
