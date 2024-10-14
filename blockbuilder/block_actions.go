package blockbuilder

import (
	"github.com/slack-go/slack"
)

// NewActionBlock creates an action block with block elements (e.g., buttons).
func NewActionBlock(elements ...slack.BlockElement) *slack.ActionBlock {
	return slack.NewActionBlock("", elements...)
}

// NewButton creates a new button element.
func NewButton(actionID, text, value string) *slack.ButtonBlockElement {
	return slack.NewButtonBlockElement(actionID, value, slack.NewTextBlockObject("plain_text", text, true, false))
}
