package models 

import "time"

type Session struct {
	SessionID string `bson:"session_id"`
	UserID string `bson:"user_id"`
	IP string `bson:"ip"`
	UserAgent string `bson:"user_agent"`
	ExpiresAt time.Time `bson:"expires_at"`
}