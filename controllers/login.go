package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/koho/dnstap-web/middleware"
	"github.com/koho/dnstap-web/models"
	"time"
)

type LoginInfo struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password"`
}

func Login(c *gin.Context) {
	req := &LoginInfo{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(400, err)
		return
	}
	u, err := models.GetUser(req.User)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	b := md5.Sum([]byte(req.Password))
	if u.Password != "" && u.Password != hex.EncodeToString(b[:]) {
		c.AbortWithStatus(403)
		return
	}
	c.Header("content-type", "text/html; charset=utf-8")
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), jwt.RegisteredClaims{
		Subject:  "admin",
		IssuedAt: jwt.NewNumericDate(time.Now()),
	})
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(middleware.SignKey)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.String(200, fmt.Sprintf("<script>localStorage.setItem('token', '%s');window.location.href='/'</script>", tokenString))
}
