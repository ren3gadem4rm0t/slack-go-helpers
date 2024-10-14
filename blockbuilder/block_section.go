package blockbuilder

import (
	"github.com/slack-go/slack"
)

// NewSectionBlock creates a new section block with optional markdown.
func NewSectionBlock(text string, markdown bool) *slack.SectionBlock {
	textObj := slack.NewTextBlockObject("mrkdwn", text, false, false)
	if !markdown {
		textObj = slack.NewTextBlockObject("plain_text", text, false, false)
	}
	return slack.NewSectionBlock(textObj, nil, nil)
}
