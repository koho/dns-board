package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koho/dns-board/middleware"
	"github.com/koho/dns-board/models"
	"net/http"
)

type LoginInfo struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password"`
}

func Login(c *gin.Context) {
	req := &LoginInfo{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	u, err := models.GetUser(req.User)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if u.Password != "" && u.Password != models.GetEncodedPassword(req.Password) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if u.Password == "" && req.Password != "" {
		models.UpdateUserPassword(u.User, req.Password)
	}
	token, err := middleware.GetToken(u)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Header("content-type", "text/html; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf("<script>localStorage.setItem('token', '%s');window.location.href='/'</script>", token))
}
