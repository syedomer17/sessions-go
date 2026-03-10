package models 

import "time"

type Session struct {
	SessionID string `bson:"session_id"`
	UserID string `bson:"user_id"`
	ExpiresAt time.Time `bson:"expires_at"`
}