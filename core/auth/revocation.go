package auth

import (
	"context"
	"time"
)

/*
Esta implementação é propositalmente neutra.
Você pode ligar em Redis, banco ou memória.
*/

type RevocationStore interface {
	IsRevoked(ctx context.Context, certFingerprint string) bool
	Revoke(ctx context.Context, certFingerprint string, ttl time.Duration) error
}
