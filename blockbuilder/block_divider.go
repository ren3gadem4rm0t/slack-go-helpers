package blockbuilder

import (
	"github.com/slack-go/slack"
)

// NewDividerBlock creates a simple divider block.
func NewDividerBlock() *slack.DividerBlock {
	return slack.NewDividerBlock()
}
