// deduper/dedupe.go

package deduper

import (
	"sync"
	"time"

	"github.com/slack-go/slack/socketmode"
)

// Dedupe handles event deduplication and caching.
type Dedupe struct {
	cache         *Cache
	evictPolicy   *EvictionPolicy
	mu            sync.Mutex
	autoEvict     bool
	evictInterval time.Duration
	stopAutoEvict chan struct{}
	wg            sync.WaitGroup
}

// Option defines a functional option for Dedupe.
type Option func(*Dedupe)

// OptionAutoEvict enables automatic eviction with a specified interval.
func OptionAutoEvict(interval time.Duration) Option {
	return func(d *Dedupe) {
		d.autoEvict = true
		d.evictInterval = interval
		d.stopAutoEvict = make(chan struct{})
	}
}

// NewDedupe initializes a new deduplication handler.
func NewDedupe(sizeLimit int, timeLimit time.Duration, countLimit int, opts ...Option) *Dedupe {
	evictPolicy := NewEvictionPolicy(sizeLimit, timeLimit, countLimit)
	d := &Dedupe{
		cache:       NewCache(evictPolicy),
		evictPolicy: evictPolicy,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(d)
	}

	// Start automatic eviction if enabled
	if d.autoEvict {
		d.wg.Add(1)
		go d.autoEvictRoutine()
	}

	return d
}

// NewDedupeWithEvictPolicy initializes a new deduplication handler with a custom eviction policy.
func NewDedupeWithEvictPolicy(evictPolicy *EvictionPolicy, opts ...Option) *Dedupe {
	d := &Dedupe{
		cache:       NewCache(evictPolicy),
		evictPolicy: evictPolicy,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(d)
	}

	// Start automatic eviction if enabled
	if d.autoEvict {
		d.wg.Add(1)
		go d.autoEvictRoutine()
	}

	return d
}

// autoEvictRoutine runs the eviction policy at specified intervals.
func (d *Dedupe) autoEvictRoutine() {
	defer d.wg.Done()
	ticker := time.NewTicker(d.evictInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.ApplyEviction(d.evictPolicy)
		case <-d.stopAutoEvict:
			return
		}
	}
}

// TriggerEviction manually triggers the eviction policy.
func (d *Dedupe) TriggerEviction() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evictPolicy.Apply(d.cache)
}

// AddEvent checks for duplicates and adds the event to the cache.
func (d *Dedupe) AddEvent(eventID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cache.Has(eventID) {
		return false // Event is a duplicate
	}

	d.cache.Add(eventID)
	return true // New event
}

// Middleware wraps a socketmode handler to add deduplication.
func (d *Dedupe) Middleware(next func(evt *socketmode.Event, client *socketmode.Client)) func(evt *socketmode.Event, client *socketmode.Client) {
	return func(evt *socketmode.Event, client *socketmode.Client) {
		// Acknowledge the event as soon as it's received if it has an envelope_id
		if evt.Request != nil {
			client.Ack(*evt.Request) // Acknowledge the event with its envelope_id
		}

		// Extract the event ID for deduplication (prefer client_msg_id if available)
		eventID, err := ExtractEventIDFromSocketMode(evt)
		if err != nil {
			client.Debugf("Failed to extract event ID: %v", err)
			return
		}

		// Add the event to the deduplication cache
		if !d.AddEvent(eventID) {
			client.Debugf("Duplicate event ignored: %s", eventID)
			return
		}

		// Proceed with the actual event handler if the event is new
		next(evt, client)
	}
}

// Size returns the current size of the deduplication cache.
func (d *Dedupe) Size() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.cache.Size()
}

// Items returns a copy of the current items in the cache.
func (d *Dedupe) Items() map[string]time.Time {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.cache.Items()
}

// ApplyEviction allows applying a custom eviction policy externally.
func (d *Dedupe) ApplyEviction(policy *EvictionPolicy) {
	d.mu.Lock()
	defer d.mu.Unlock()
	policy.Apply(d.cache)
}

// StopAutoEviction stops the automatic eviction goroutine.
func (d *Dedupe) StopAutoEviction() {
	if d.autoEvict {
		close(d.stopAutoEvict)
		d.wg.Wait()
	}
}
