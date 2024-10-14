package blockbuilder

import (
	"github.com/slack-go/slack"
)

// NewContextBlock creates a context block with mixed elements (text/images).
func NewContextBlock(elements ...slack.MixedElement) *slack.ContextBlock {
	return slack.NewContextBlock("", elements...)
}
