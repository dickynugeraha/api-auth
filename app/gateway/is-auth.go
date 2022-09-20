package gateway

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func IsAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error" : "No header authorization",
			})
			return
		}
		

		token := 	c.Query("token")
		
		bearerToken := c.Request.Header.Get("Authorization")
		if len(strings.Split(bearerToken, " ")) == 2 {
			token = strings.Split(bearerToken, " ")[1]
		}

		token1, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			if _, ok := token1.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signin method %v", token1.Header["alg"])
			}
			return secret, nil
		})
	}
}