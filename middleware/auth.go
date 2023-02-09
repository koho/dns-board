package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/koho/dns-board/db"
	"github.com/koho/dns-board/models"
	"net/http"
	"strings"
	"time"
)

var signKey []byte

func init() {
	db.OnStartup(func(option db.Option) error {
		sk, err := models.GetMeta("sk")
		if err != nil {
			k := make([]byte, 16)
			if _, err = rand.Read(k); err != nil {
				return err
			}
			if err = models.SetMeta("sk", base64.StdEncoding.EncodeToString(k)); err != nil {
				return err
			}
			signKey = k
			return nil
		}
		k, err := base64.StdEncoding.DecodeString(sk)
		if err != nil {
			return err
		}
		signKey = k
		return nil
	})
}

func GetToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), jwt.RegisteredClaims{
		Subject:  user.User,
		IssuedAt: jwt.NewNumericDate(time.Now()),
	})
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer"))
		if token != "" {
			if _, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
				return signKey, nil
			}); err == nil {
				c.Next()
				return
			}
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
