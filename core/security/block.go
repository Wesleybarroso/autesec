package security

import (
	"sync"
	"time"
)

type BlockList struct {
	mu    sync.Mutex
	items map[string]time.Time
}

func NewBlockList() *BlockList {
	return &BlockList{
		items: make(map[string]time.Time),
	}
}

func (b *BlockList) Block(key string, duration time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.items[key] = time.Now().Add(duration)
}

func (b *BlockList) IsBlocked(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	exp, ok := b.items[key]
	if !ok {
		return false
	}

	if time.Now().After(exp) {
		delete(b.items, key)
		return false
	}

	return true
}
