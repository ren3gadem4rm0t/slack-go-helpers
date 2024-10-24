// examples/deduper/helpers/main.go

package main

import (
	"fmt"

	"github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	// Assume you have a socketmode.Event named evt
	var evt socketmode.Event
	// Normally, evt would be populated with real event data

	// Extract the event ID
	eventID, err := deduper.ExtractEventIDFromSocketMode(&evt)
	if err != nil {
		fmt.Printf("Error extracting event ID: %v\n", err)
	} else {
		fmt.Printf("Extracted Event ID: %s\n", eventID)
	}

	// Alternatively, extract the envelope ID
	envelopeID, err := deduper.ExtractEnvelopeIDFromSocketMode(&evt)
	if err != nil {
		fmt.Printf("Error extracting envelope ID: %v\n", err)
	} else {
		fmt.Printf("Extracted Envelope ID: %s\n", envelopeID)
	}
}
