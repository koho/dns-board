package controllers

import (
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/gin-gonic/gin"
	"github.com/koho/dns-board/db"
	"github.com/koho/dns-board/models"
	"net/http"
)

func GetCacheHit(c *gin.Context) {
	m := db.GetDB().Model(&models.Message{}).Select("strftime('%Y-%m-%d %H:%M', time, 'localtime') as t, round(sum(iif(duration = 0, 1, 0)) / cast(count(*) as real), 2) as value").
		Where("type = ?", dnstap.Message_CLIENT_RESPONSE)
	m = models.AddTimeClause(m, c.DefaultQuery("hour", "3"))
	var hits []TimeSeries
	if err := m.Group("t").Find(&hits).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, hits)
}
