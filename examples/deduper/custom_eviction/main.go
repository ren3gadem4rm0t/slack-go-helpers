// examples/deduper/custom_eviction/main.go

package main

import (
	"fmt"
	"time"

	"github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
)

// CustomEvictionPolicy defines a custom eviction policy
type CustomEvictionPolicy struct {
	*deduper.EvictionPolicy
}

// Apply overrides the default Apply method to implement custom logic
func (e *CustomEvictionPolicy) Apply(c *deduper.Cache) {
	// Custom rule: Evict events that start with "temp"
	for id := range c.Items() {
		if len(id) >= 4 && id[:4] == "temp" {
			c.Evict(id)
			fmt.Printf("Evicted event based on custom rule: %s\n", id)
		}
	}

	// Call the original eviction logic
	e.EvictionPolicy.Apply(c)
}

func main() {
	// Initialize the custom eviction policy
	baseEvictionPolicy := deduper.NewEvictionPolicy(1000, 10*time.Minute, 1000)
	customEvictionPolicy := &CustomEvictionPolicy{baseEvictionPolicy}

	// Initialize the deduplication handler with the custom eviction policy
	dedupeHandler := deduper.NewDedupeWithEvictPolicy(customEvictionPolicy.EvictionPolicy) // Corrected line

	// Add events
	events := []string{"temp_event1", "event2", "temp_event3", "event4"}

	for _, id := range events {
		dedupeHandler.AddEvent(id)
	}

	fmt.Printf("Cache size before eviction: %d\n", dedupeHandler.Size())

	// Manually trigger eviction using the ApplyEviction method
	dedupeHandler.ApplyEviction(customEvictionPolicy.EvictionPolicy)

	fmt.Printf("Cache size after eviction: %d\n", dedupeHandler.Size())

	// Check remaining events
	for id := range dedupeHandler.Items() {
		fmt.Printf("Remaining event: %s\n", id)
	}
}
