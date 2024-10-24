// deduper/cache.go

package deduper

import (
	"sync"
	"time"
)

// Cache stores event IDs with safe concurrent access.
type Cache struct {
	items       map[string]time.Time
	mu          sync.RWMutex
	evictPolicy *EvictionPolicy
}

// NewCache initializes a new cache.
func NewCache(evictPolicy *EvictionPolicy) *Cache {
	return &Cache{
		items:       make(map[string]time.Time),
		evictPolicy: evictPolicy,
	}
}

// Has checks if the eventID is already in the cache.
func (c *Cache) Has(eventID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.items[eventID]
	return exists
}

// Add inserts a new eventID into the cache and applies eviction if needed.
func (c *Cache) Add(eventID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[eventID] = time.Now()
	c.evictPolicy.Apply(c)
}

// Size returns the current size of the cache.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Evict removes an eventID from the cache.
func (c *Cache) Evict(eventID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, eventID)
}

// Items returns a copy of the current items in the cache.
func (c *Cache) Items() map[string]time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	copiedItems := make(map[string]time.Time, len(c.items))
	for k, v := range c.items {
		copiedItems[k] = v
	}
	return copiedItems
}
