package middleware

import (
	"net/http"
	"session-demo/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")

		if err != nil || sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: No session cookie",
			})
			return
		}
		session , err := services.FindSession(sessionID)

		if err != nil || session.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "session expires",
			})
			return
		}

		if session.IP != c.ClientIP() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: IP address mismatch",
			})
			return
		}

		if session.UserAgent != c.Request.UserAgent() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: User agent mismatch",
			})
			return
		}

		c.Set("user_id", session.UserID)

		c.Next()
	}
}