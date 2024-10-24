// examples/deduper/comprehensive/main.go

package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	// Initialize the deduplication handler with size limit, time limit, and count limit
	dedupeHandler := deduper.NewDedupe(1000, 10*time.Minute, 1000)

	// Create a new Slack client with the App Level Token
	slackClient := slack.New(
		os.Getenv("SLACK_BOT_TOKEN"),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")),
	)

	// Create a new Socket Mode client
	client := socketmode.New(
		slackClient,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	// Define your event handler
	eventHandler := func(evt *socketmode.Event, client *socketmode.Client) {
		switch evt.Type {
		case socketmode.EventTypeEventsAPI:
			var msg slackevents.EventsAPIEvent
			data, ok := evt.Data.(json.RawMessage)
			if !ok {
				log.Printf("Failed to assert evt.Data to json.RawMessage")
				return
			}
			err := json.Unmarshal(data, &msg)
			if err != nil {
				log.Printf("Failed to unmarshal EventsAPIEvent: %v", err)
				return
			}

			// Acknowledge the event
			client.Ack(*evt.Request)

			// Handle the event via middleware
			wrappedHandler := dedupeHandler.Middleware(func(evt *socketmode.Event, client *socketmode.Client) {
				// Further event processing based on msg.Type
				switch msg.Type {
				case slackevents.CallbackEvent:
					innerEvent := msg.InnerEvent
					log.Printf("Received inner event: %v", innerEvent)
				default:
					log.Printf("Unhandled EventsAPI event type: %s", msg.Type)
				}
			})
			wrappedHandler(evt, client)

		case socketmode.EventTypeInteractive:
			wrappedHandler := dedupeHandler.Middleware(func(evt *socketmode.Event, client *socketmode.Client) {
				var payload *slack.InteractionCallback
				data, ok := evt.Data.(json.RawMessage)
				if !ok {
					client.Debugf("Failed to assert evt.Data to json.RawMessage")
					return
				}
				err := json.Unmarshal(data, &payload)
				if err != nil {
					client.Debugf("Could not unmarshal interaction callback: %v", err)
					return
				}

				log.Printf("Button clicked: %s", payload.ActionCallback.BlockActions[0].Value)
			})
			wrappedHandler(evt, client)

		case socketmode.EventTypeSlashCommand:
			wrappedHandler := dedupeHandler.Middleware(func(evt *socketmode.Event, client *socketmode.Client) {
				var cmd slack.SlashCommand
				data, ok := evt.Data.(json.RawMessage)
				if !ok {
					client.Debugf("Failed to assert evt.Data to json.RawMessage")
					return
				}
				err := json.Unmarshal(data, &cmd)
				if err != nil {
					client.Debugf("Could not unmarshal slash command: %v", err)
					return
				}

				log.Printf("Slash command received: %s", cmd.Command)
			})
			wrappedHandler(evt, client)

		default:
			// Ignore other event types
		}
	}

	// Listen for events
	go func() {
		for evt := range client.Events {
			eventHandler(&evt, client)
		}
	}()

	log.Println("Starting Slack Socket Mode client...")
	// Run the client
	err := client.Run()
	if err != nil {
		log.Fatalf("Failed to run client: %v", err)
	}
}
