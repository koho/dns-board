package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/koho/dns-board/models"
	"net/http"
)

func SearchRecord(c *gin.Context) {
	kw := c.DefaultQuery("kw", "")
	if kw == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	records, err := models.GetMessagesByKeyword(kw)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, records)
}
