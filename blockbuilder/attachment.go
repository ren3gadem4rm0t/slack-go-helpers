package blockbuilder

import (
	"github.com/slack-go/slack"
)

// AttachmentBuilder is used to create Slack message attachments with optional color and blocks.
type AttachmentBuilder struct {
	attachment slack.Attachment
}

// NewAttachmentBuilder initializes a new attachment builder with an optional color.
func NewAttachmentBuilder(color string) *AttachmentBuilder {
	return &AttachmentBuilder{
		attachment: slack.Attachment{
			Color: color, // Set the color bar for the attachment.
		},
	}
}

// AddSection adds a section block to the attachment.
func (a *AttachmentBuilder) AddSection(text string, markdown bool) *AttachmentBuilder {
	section := NewSectionBlock(text, markdown)
	a.attachment.Blocks.BlockSet = append(a.attachment.Blocks.BlockSet, section)
	return a
}

// AddActions adds an action block with buttons to the attachment.
func (a *AttachmentBuilder) AddActions(elements ...slack.BlockElement) *AttachmentBuilder {
	actions := NewActionBlock(elements...)
	a.attachment.Blocks.BlockSet = append(a.attachment.Blocks.BlockSet, actions)
	return a
}

// AddContext adds a context block to the attachment.
func (a *AttachmentBuilder) AddContext(elements ...slack.MixedElement) *AttachmentBuilder {
	context := NewContextBlock(elements...)
	a.attachment.Blocks.BlockSet = append(a.attachment.Blocks.BlockSet, context)
	return a
}

// AddImage adds an image block to the attachment.
func (a *AttachmentBuilder) AddImage(imageURL, altText string) *AttachmentBuilder {
	image := NewImageBlock(imageURL, altText)
	a.attachment.Blocks.BlockSet = append(a.attachment.Blocks.BlockSet, image)
	return a
}

// AddDivider adds a divider block to the attachment.
func (a *AttachmentBuilder) AddDivider() *AttachmentBuilder {
	divider := NewDividerBlock()
	a.attachment.Blocks.BlockSet = append(a.attachment.Blocks.BlockSet, divider)
	return a
}

// Build returns the constructed attachment.
func (a *AttachmentBuilder) Build() slack.Attachment {
	return a.attachment
}
