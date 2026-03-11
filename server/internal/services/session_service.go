package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"session-demo/internal/config"
	"session-demo/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GenerateSessionID() string {
	bytes := make([]byte, 32)

	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}

func CreateSession(userID string, ip string, userAgent string) (string, error) {
	sessionID := GenerateSessionID()

	session := models.Session{
		SessionID: sessionID,
		UserID:    userID,
		IP:        ip,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	collection := config.DB.Collection("sessions")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, session)

	if err != nil {
		return "", err
	}

	config.RedisClient.Set(
		ctx,
		sessionID,
		userID,
		24*time.Hour,
	)

	return sessionID, nil
}

func FindSession(sessionID string) (*models.Session, error) {
	collections := config.DB.Collection("sessions")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session models.Session

	err := collections.FindOne(ctx, bson.M{
		"session_id": sessionID,
	}).Decode(&session)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func DeleteSession(sessionID string) error {
	collections := config.DB.Collection("sessions")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collections.DeleteOne(ctx, bson.M{
		"session_id": sessionID,
	})

	if err != nil {
		return err
	}

	return nil
}

func RotateSession(oldSessionID string, userID string) (string, error) {
	err := DeleteSession(oldSessionID)

	if err != nil {
		return "", err
	}

	newSessionID, err := CreateSession(userID, "", "")

	if err != nil {
		return "", err
	}

	return newSessionID, nil
}
