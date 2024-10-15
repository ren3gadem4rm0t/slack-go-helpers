package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ren3gadem4rm0t/slack-go-helpers/blockbuilder"
	"github.com/slack-go/slack"
)

func main() {
	// Create a BlockBuilder instance to build a custom block.
	customBlock := blockbuilder.NewBlockBuilder().
		AddSection("This is a custom block added using AddBlock!", true).
		AddDivider().
		AddImage("https://image.url/example.png", "Example Image").
		Build()

	// Create an attachment and add a combination of custom blocks and predefined blocks.
	attachmentBuilder := blockbuilder.NewAttachmentBuilder(blockbuilder.ColorWarning).
		AddSection("This is a warning message", true).
		AddDivider().
		AddBlock(customBlock[0]). // Adding a single custom block
		AddActions(
			blockbuilder.NewButton("ack_btn", "Acknowledge", "warn_value_1"),
		).
		AddContext(slack.NewTextBlockObject("plain_text", "Proceed with caution.", false, false))

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
