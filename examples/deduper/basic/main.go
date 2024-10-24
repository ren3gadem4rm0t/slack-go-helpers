// examples/deduper/basic/main.go

package main

import (
	"fmt"
	"time"

	"github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
)

func main() {
	// Initialize the deduplication handler with size limit, time limit, and count limit
	dedupeHandler := deduper.NewDedupe(1000, 5*time.Minute, 500)

	eventIDs := []string{"event1", "event2", "event3", "event1", "event2"}

	for _, id := range eventIDs {
		if dedupeHandler.AddEvent(id) {
			fmt.Printf("Processed event: %s\n", id)
		} else {
			fmt.Printf("Duplicate event ignored: %s\n", id)
		}
	}

	fmt.Printf("Current cache size: %d\n", dedupeHandler.Size())
}
