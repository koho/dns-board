package middleware

import (
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"os"
	"strings"
)

var SignKey []byte

func init() {
	sk, err := os.ReadFile("sign.key")
	if err != nil {
		sk = make([]byte, 16)
		if _, err = rand.Read(sk); err != nil {
			log.Fatal(err)
		}
		if err = os.WriteFile("sign.key", sk, 0666); err != nil {
			log.Fatal(err)
		}
	}
	SignKey = sk
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer"))
		if token != "" {
			if _, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
				return SignKey, nil
			}); err == nil {
				c.Next()
				return
			}
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
