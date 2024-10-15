package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ren3gadem4rm0t/slack-go-helpers/blockbuilder"
	"github.com/slack-go/slack"
)

func main() {
	blockBuilder := blockbuilder.NewBlockBuilder().
		AddSection("Welcome to the app!", true).
		AddDivider().
		AddActions(
			blockbuilder.NewButton("btn_1", "Click Me", "value_1"),
			blockbuilder.NewButton("btn_2", "Cancel", "value_2"),
		).
		AddContext(slack.NewTextBlockObject("plain_text", "Powered by Slack-Go BlockBuilder", false, false)).
		AddImage("https://image.url/example.png", "Example Image")

	// Build the blocks and attachments separately.
	blocks := blockBuilder.Build()

	// Marshal the blocks to JSON.
	blocksJSON, err := json.MarshalIndent(blocks, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling blocks: %v", err)
	}
	// Print the JSON representation of blocks and attachment.
	fmt.Println("Blocks JSON:")
	fmt.Println(string(blocksJSON))

	// Expected output:
	// Blocks JSON:
	// [
	// 	{
	// 		"type": "section",
	// 		"text": {
	// 		"type": "mrkdwn",
	// 		"text": "Welcome to the app!"
	// 		}
	// 	},
	// 	{
	// 		"type": "divider"
	// 	},
	// 	{
	// 		"type": "actions",
	// 		"elements": [
	// 		{
	// 			"type": "button",
	// 			"text": {
	// 			"type": "plain_text",
	// 			"text": "Click Me",
	// 			"emoji": true
	// 			},
	// 			"action_id": "btn_1",
	// 			"value": "value_1"
	// 		},
	// 		{
	// 			"type": "button",
	// 			"text": {
	// 			"type": "plain_text",
	// 			"text": "Cancel",
	// 			"emoji": true
	// 			},
	// 			"action_id": "btn_2",
	// 			"value": "value_2"
	// 		}
	// 		]
	// 	},
	// 	{
	// 		"type": "context",
	// 		"elements": [
	// 		{
	// 			"type": "plain_text",
	// 			"text": "Powered by Slack-Go BlockBuilder"
	// 		}
	// 		]
	// 	},
	// 	{
	// 		"type": "image",
	// 		"image_url": "https://image.url/example.png",
	// 		"alt_text": "Example Image"
	// 	}
	// ]
}
