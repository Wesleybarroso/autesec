package security

import (
	"sync"
	"time"
)

type fingerprintStat struct {
	count int
	reset time.Time
}

type AnomalyDetector struct {
	mu        sync.Mutex
	threshold int
	window    time.Duration
	stats     map[string]*fingerprintStat
}

func NewAnomalyDetector(threshold int, window time.Duration) *AnomalyDetector {
	return &AnomalyDetector{
		threshold: threshold,
		window:    window,
		stats:     make(map[string]*fingerprintStat),
	}
}

func (a *AnomalyDetector) Register(fingerprint string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	s, ok := a.stats[fingerprint]
	if !ok || time.Now().After(s.reset) {
		a.stats[fingerprint] = &fingerprintStat{
			count: 1,
			reset: time.Now().Add(a.window),
		}
		return true
	}

	s.count++
	return s.count <= a.threshold
}
