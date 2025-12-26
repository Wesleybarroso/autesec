package session

import "context"

type Store interface {
	Save(ctx context.Context, s *Session) error
	Get(ctx context.Context, sessionID string) (*Session, error)
	Delete(ctx context.Context, sessionID string) error
	UpdateLastSeen(ctx context.Context, sessionID string) error
}
