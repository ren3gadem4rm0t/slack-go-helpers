package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ren3gadem4rm0t/slack-go-helpers/blockbuilder"
	"github.com/slack-go/slack"
)

func main() {
	// Create a BlockBuilder instance to build custom blocks.
	blockBuilder := blockbuilder.NewBlockBuilder().
		AddSection("Custom block added using BlockBuilder!", true).
		AddImage("https://image.url/example.png", "Example Image")

	// Create an attachment and add both custom blocks and blocks through helper methods.
	attachmentBuilder := blockbuilder.NewAttachmentBuilder(blockbuilder.ColorWarning).
		AddSection("This is a warning message", true).
		AddDivider().
		AddBlocksFromBuilder(blockBuilder). // Add blocks from the BlockBuilder
		AddActions(
			blockbuilder.NewButton("ack_btn", "Acknowledge", "warn_value_1"),
		).
		AddContext(slack.NewTextBlockObject("plain_text", "Please proceed with caution.", false, false))

	// Build the attachment.
	attachment := attachmentBuilder.Build()

	// Marshal the attachment to JSON for display.
	attachmentJSON, err := json.MarshalIndent(attachment, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling attachment: %v", err)
	}

	// Print the JSON representation of the attachment.
	fmt.Println("Attachment JSON:")
	fmt.Println(string(attachmentJSON))
}
