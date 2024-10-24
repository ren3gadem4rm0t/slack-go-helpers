package deduper

import (
	"time"
)

// EvictionPolicy defines rules for removing stale entries.
type EvictionPolicy struct {
	sizeLimit  int
	timeLimit  time.Duration
	countLimit int
}

// NewEvictionPolicy creates a new eviction policy.
func NewEvictionPolicy(sizeLimit int, timeLimit time.Duration, countLimit int) *EvictionPolicy {
	return &EvictionPolicy{
		sizeLimit:  sizeLimit,
		timeLimit:  timeLimit,
		countLimit: countLimit,
	}
}

// Apply evicts items that violate the policy.
func (e *EvictionPolicy) Apply(c *Cache) {
	now := time.Now()

	// Evict based on time limit
	for id, t := range c.items {
		if now.Sub(t) > e.timeLimit {
			delete(c.items, id)
		}
	}

	// Evict based on count limit
	if len(c.items) > e.countLimit {
		for id := range c.items {
			delete(c.items, id)
			if len(c.items) <= e.countLimit {
				break
			}
		}
	}
}
