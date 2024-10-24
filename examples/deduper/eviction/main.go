// examples/deduper/eviction/main.go

package main

import (
	"fmt"
	"time"

	"github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
)

func main() {
	// Create an eviction policy with a size limit of 3, time limit of 10 seconds, and count limit of 3
	evictionPolicy := deduper.NewEvictionPolicy(3, 10*time.Second, 3)
	cache := deduper.NewCache(evictionPolicy)

	// Add some events
	cache.Add("event1")
	cache.Add("event2")
	cache.Add("event3")

	fmt.Printf("Initial cache size: %d\n", cache.Size())

	// Add a duplicate event
	cache.Add("event2")

	fmt.Printf("Cache size after adding duplicate: %d\n", cache.Size())

	// Add a new event to trigger eviction based on size limit
	cache.Add("event4")

	fmt.Printf("Cache size after adding event4: %d\n", cache.Size())

	// Wait for eviction based on time limit
	time.Sleep(11 * time.Second)
	cache.Add("event5")

	fmt.Printf("Cache size after time-based eviction: %d\n", cache.Size())
}
