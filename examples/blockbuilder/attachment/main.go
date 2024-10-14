package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ren3gadem4rm0t/slack-go-helpers/blockbuilder"
	"github.com/slack-go/slack"
)

func main() {
	// Create a block builder for the main message.
	blockBuilder := blockbuilder.NewBlockBuilder().
		AddSection("Welcome to the app!", true).
		AddDivider().
		AddActions(
			blockbuilder.NewButton("btn_1", "Click Me", "value_1"),
			blockbuilder.NewButton("btn_2", "Cancel", "value_2"),
		).
		AddContext(slack.NewTextBlockObject("plain_text", "Powered by Slack", false, false))

	// Create an attachment with a color bar and some content.
	attachmentBuilder := blockbuilder.NewAttachmentBuilder(blockbuilder.ColorWarning).
		AddSection("This is a warning message", true).
		AddDivider().
		AddActions(
			blockbuilder.NewButton("warn_btn_1", "Acknowledge", "warn_value_1"),
		).
		AddContext(slack.NewTextBlockObject("plain_text", "Please be cautious", false, false))

	// Build the blocks and attachments separately.
	blocks := blockBuilder.Build()
	attachment := attachmentBuilder.Build()

	// Marshal the blocks to JSON.
	blocksJSON, err := json.MarshalIndent(blocks, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling blocks: %v", err)
	}

	// Marshal the attachment to JSON.
	attachmentJSON, err := json.MarshalIndent([]slack.Attachment{attachment}, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling attachments: %v", err)
	}

	// Print the JSON representation of blocks and attachment.
	fmt.Println("Blocks JSON:")
	fmt.Println(string(blocksJSON))

	fmt.Println("\nAttachment JSON:")
	fmt.Println(string(attachmentJSON))
}
