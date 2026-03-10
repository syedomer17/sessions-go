package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohugoio/hugo/tpl/collections"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func login(c *gin.Context){
	sessionID := uuid.New().String()

	session := bson.M{
		"session_id": sessionID,
		"user_id": "12345",
		"expires_at": time.Now().Add(24 * time.Hour),
	}

	collections.Sessions.InsertOne(context.TODO(), session)

	c.Cookie(
		"session_id",
		sessionID,
		3600*24,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "logged in",
	})
}