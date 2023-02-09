package controllers

import (
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/gin-gonic/gin"
	"github.com/koho/dns-board/db"
	"github.com/koho/dns-board/models"
	"net/http"
)

type QueryIPStat struct {
	IP    string `json:"ip"`
	Count int64  `json:"count"`
}

func GetQueryIPStat(c *gin.Context) {
	m := db.GetDB().Model(&models.Message{}).Select("ip, count(*) as count").
		Where("type = ? and answer <> ''", dnstap.Message_CLIENT_RESPONSE)
	var ipStat []QueryIPStat
	if err := models.AddTimeClause(m, c.DefaultQuery("hour", "3")).
		Group("ip").Order("count desc").Find(&ipStat).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ipStat)
}

type QueryTypeStat struct {
	QType string `json:"type"`
	Count int64  `json:"value"`
}

func GetQueryTypeStat(c *gin.Context) {
	m := db.GetDB().Model(&models.Message{}).Select("q_type, count(*) as count").
		Where("type = ? and answer <> ''", dnstap.Message_CLIENT_RESPONSE)
	var typeStat []QueryTypeStat
	if err := models.AddTimeClause(m, c.DefaultQuery("hour", "3")).
		Group("q_type").Order("count desc").Find(&typeStat).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, typeStat)
}

type QueryCountStat struct {
	T     string `json:"time"`
	Value int64  `json:"value"`
}

func GetQueryCountStat(c *gin.Context) {
	m := db.GetDB().Model(&models.Message{}).Select("strftime('%Y-%m-%d %H:%M', time, 'localtime') as t, count(*) as value").
		Where("type = ?", dnstap.Message_CLIENT_RESPONSE)
	var countStat []QueryCountStat
	if err := models.AddTimeClause(m, c.DefaultQuery("hour", "3")).Group("t").Find(&countStat).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, countStat)
}
