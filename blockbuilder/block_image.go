package blockbuilder

import (
	"github.com/slack-go/slack"
)

// NewImageBlock creates an image block with a URL and alt text.
func NewImageBlock(imageURL, altText string) *slack.ImageBlock {
	return slack.NewImageBlock(imageURL, altText, "", nil)
}
