package session

import "time"

type Session struct {
	SessionID       string
	UserID          string
	CertFingerprint string
	CreatedAt       time.Time
	LastSeen        time.Time
}
