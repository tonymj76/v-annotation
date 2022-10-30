package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const (
	authorizationHeader = "Authorization"
	apiKeyHeader        = "X-API-Key"
)

// AuthAPIKey key auth middleware
func AuthAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get(apiKeyHeader)

		secret := os.Getenv("SESSION_SECRET")
		if secret == "" {
			logrus.Println("failed to get secret")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		logrus.Println("comparing secret with provided key", secret, key)

		if secret != key {
			logrus.Println("key doesnt match!")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		logrus.Println("no error during check")
		c.Next()
	}
}
