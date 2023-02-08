package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/koho/dnstap-web/models"
)

func SearchRecord(c *gin.Context) {
	kw := c.DefaultQuery("kw", "")
	if kw == "" {
		c.AbortWithStatus(400)
		return
	}
	records, err := models.GetMessagesByKeyword(kw)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, records)
}
