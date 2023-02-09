package controllers

import (
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/gin-gonic/gin"
	"github.com/koho/dns-board/db"
	"github.com/koho/dns-board/models"
	"net/http"
)

func GetResponseSize(c *gin.Context) {
	m := db.GetDB().Model(&models.Message{}).Select("strftime('%Y-%m-%d %H:%M', time, 'localtime') as t, round(avg(size), 0) as value").
		Where("type = ?", dnstap.Message_CLIENT_RESPONSE)
	var sizeStat []TimeSeries
	if err := models.AddTimeClause(m, c.DefaultQuery("hour", "3")).Group("t").Find(&sizeStat).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, sizeStat)
}
