package utils

import "github.com/gin-gonic/gin"

func SetSessionCookie(sessionID string, c *gin.Context){
	c.SetCookie(
		"session_id",
		sessionID,
		3600 * 24,
		"/",
		"localhost",
		false,
		true,
	)
}

func ClearSessionCookie(c *gin.Context){
	c.SetCookie(
		"session_id",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)
}