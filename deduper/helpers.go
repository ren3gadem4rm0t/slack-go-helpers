package deduper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// ExtractEventIDFromSocketMode extracts a stable event ID for deduplication.
// Uses client_msg_id for message events, otherwise hashes key event fields.
func ExtractEventIDFromSocketMode(evt *socketmode.Event) (string, error) {
	// Check if the event is an EventsAPI event (events sent via the Events API)
	event, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		return "", errors.New("invalid event in socketmode event")
	}

	// Prefer client_msg_id for message events
	if messageEvent, ok := event.InnerEvent.Data.(*slackevents.MessageEvent); ok && messageEvent.ClientMsgID != "" {
		return messageEvent.ClientMsgID, nil
	}

	// For other events, create an MD5 hash of key parts of the event
	stableID, err := generateEventHash(event)
	if err != nil {
		return "", err
	}

	return stableID, nil
}

// generateEventHash creates an MD5 hash from stable parts of the event for deduplication.
func generateEventHash(event slackevents.EventsAPIEvent) (string, error) {
	// Extract core parts of the event to hash (these fields are generally consistent across retries)
	// Include more fields like event type, team ID, user ID, timestamp, and channel ID if available
	var userID, eventTs, channelID, customID string

	// Handle different inner events to extract stable IDs
	switch ev := event.InnerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		userID = ev.User
		eventTs = ev.TimeStamp
		channelID = ev.Channel
		customID = ev.ClientMsgID
	case *slackevents.ReactionAddedEvent:
		userID = ev.User
		eventTs = ev.EventTimestamp
		channelID = ev.Item.Channel
		customID = ev.Reaction
	case *slackevents.AppMentionEvent:
		userID = ev.User
		eventTs = ev.TimeStamp
		channelID = ev.Channel
		customID = ev.Text
	case *slackevents.ReactionRemovedEvent:
		userID = ev.User
		eventTs = ev.EventTimestamp
		channelID = ev.Item.Channel
		customID = ev.Reaction
	default:
		return "", errors.New("unsupported event type")
	}

	// Fallback if no eventTs was found
	if eventTs == "" {
		eventTs = "no-timestamp"
	}

	// Create a string containing the relevant information
	coreData := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", event.Type, event.TeamID, userID, eventTs, channelID, event.InnerEvent.Type, customID)

	// Generate MD5 hash
	hash := md5.New()
	_, err := hash.Write([]byte(coreData))
	if err != nil {
		return "", err
	}

	// Return the hash as a hexadecimal string
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ExtractEventIDFromSocketMode extracts the envelope ID from a socketmode event.
func ExtractEnvelopeIDFromSocketMode(evt *socketmode.Event) (string, error) {
	var evId string
	if evt.Request != nil {
		if evt.Request.EnvelopeID == "" {
			return "", errors.New("event does not have an envelope_id")
		}
		evId = evt.Request.EnvelopeID
	}

	// Return the envelope_id for deduplication
	return evId, nil
}
