// examples/deduper/testing/main.go

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
)

func main() {
	// Initialize the deduplication handler with size limit, time limit, and count limit
	// Enable automatic eviction with an interval of 20 seconds
	dedupeHandler := deduper.NewDedupe(1000, 20*time.Second, 1000, deduper.OptionAutoEvict(10*time.Second))
	defer dedupeHandler.StopAutoEviction() // Ensure the eviction goroutine is stopped when main exits

	var wg sync.WaitGroup

	// Include some "temp_*" events to test manual eviction
	eventIDs := []string{
		"event1",
		"event2",
		"temp_event3", // Should be evicted manually
		"event1",      // Duplicate
		"event2",      // Duplicate
		"event4",
		"temp_event5", // Should be evicted manually
		"event6",
	}

	for _, id := range eventIDs {
		wg.Add(1)
		go func(eventID string) {
			defer wg.Done()
			if eventID == "event6" {
				// Wait for 18 seconds before processing event6 to allow some events to expire
				time.Sleep(18 * time.Second)
			}
			if dedupeHandler.AddEvent(eventID) {
				fmt.Printf("Processed event: %s\n", eventID)
			} else {
				fmt.Printf("Duplicate event ignored: %s\n", eventID)
			}
		}(id)
	}

	wg.Wait()

	fmt.Printf("Final cache size: %d\n", dedupeHandler.Size())

	// Manually trigger eviction to remove "temp_*" events
	fmt.Println("Manually triggering eviction...")
	dedupeHandler.TriggerEviction()

	// After manual eviction, check cache size and contents
	fmt.Printf("Cache size after manual eviction: %d\n", dedupeHandler.Size())
	fmt.Println("Remaining events after manual eviction:")
	for id := range dedupeHandler.Items() {
		fmt.Printf("- %s\n", id)
	}

	// Wait for automatic eviction to occur
	fmt.Println("Waiting for automatic eviction to occur...")
	time.Sleep(15 * time.Second)

	// Final cache size after automatic eviction
	fmt.Printf("Cache size after automatic eviction: %d\n", dedupeHandler.Size())
	fmt.Println("Remaining events after automatic eviction:")
	for id := range dedupeHandler.Items() {
		fmt.Printf("- %s\n", id)
	}
}

// Expected Output:
// $ go run examples/deduper/testing/main.go
// Processed event: temp_event3
// Processed event: temp_event5
// Processed event: event2
// Processed event: event4
// Duplicate event ignored: event1
// Duplicate event ignored: event2
// Processed event: event1
// Processed event: event6
// Final cache size: 6
// Manually triggering eviction...
// Cache size after manual eviction: 6
// Remaining events after manual eviction:
// - temp_event5
// - event6
// - temp_event3
// - event1
// - event4
// - event2
// Waiting for automatic eviction to occur...
// Cache size after automatic eviction: 1
// Remaining events after automatic eviction:
// - event6
