package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/koho/dns-board/db"
	"github.com/koho/dns-board/models"
	"net/http"
)

func GetRawData(c *gin.Context) {
	var msg []models.Message
	if err := models.AddTimeClause(db.GetDB(), c.DefaultQuery("hour", "3")).Find(&msg).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, msg)
}
