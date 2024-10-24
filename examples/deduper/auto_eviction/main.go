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
	eventIDs := []string{"event1", "event2", "event3", "event1", "event2", "event4", "event5", "event6"}

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

	// Manually trigger eviction
	fmt.Println("Manually triggering eviction...")
	dedupeHandler.TriggerEviction()

	// After eviction, check cache size
	fmt.Printf("Cache size after manual eviction: %d\n", dedupeHandler.Size())

	// Wait for automatic eviction to occur
	fmt.Println("Waiting for automatic eviction to occur...")
	time.Sleep(15 * time.Second)

	// Final cache size after automatic eviction
	fmt.Printf("Cache size after automatic eviction: %d\n", dedupeHandler.Size())
}
