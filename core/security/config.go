package security

import "time"

type Config struct {
	MaxConnPerIP      int
	Window            time.Duration
	AnomalyThreshold  int
	BlockDuration     time.Duration
}
