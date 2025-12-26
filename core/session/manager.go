package session

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Manager struct {
	store Store
}

func NewManager(store Store) *Manager {
	return &Manager{store: store}
}

func (m *Manager) Create(ctx context.Context, userID, certFP string) (*Session, error) {
	s := &Session{
		SessionID:       uuid.NewString(),
		UserID:          userID,
		CertFingerprint: certFP,
		CreatedAt:       time.Now(),
		LastSeen:        time.Now(),
	}

	if err := m.store.Save(ctx, s); err != nil {
		return nil, err
	}

	return s, nil
}

func (m *Manager) Resume(ctx context.Context, sessionID, certFP string) (*Session, error) {
	s, err := m.store.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if s.CertFingerprint != certFP {
		return nil, errors.New("session certificate mismatch")
	}

	s.LastSeen = time.Now()

	if err := m.store.Save(ctx, s); err != nil {
		return nil, err
	}

	return s, nil
}

func (m *Manager) Close(ctx context.Context, sessionID string) error {
	return m.store.Delete(ctx, sessionID)
}
