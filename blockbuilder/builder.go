package blockbuilder

import (
	"github.com/slack-go/slack"
)

// BlockBuilder builds a message composed of blocks.
type BlockBuilder struct {
	blocks []slack.Block
}

// NewBlockBuilder initializes an empty BlockBuilder.
func NewBlockBuilder() *BlockBuilder {
	return &BlockBuilder{
		blocks: []slack.Block{},
	}
}

// AddSection adds a section block with text to the builder.
func (b *BlockBuilder) AddSection(text string, markdown bool) *BlockBuilder {
	section := NewSectionBlock(text, markdown)
	b.blocks = append(b.blocks, section)
	return b
}

// AddActions adds an action block with buttons to the builder.
func (b *BlockBuilder) AddActions(elements ...slack.BlockElement) *BlockBuilder {
	actions := NewActionBlock(elements...)
	b.blocks = append(b.blocks, actions)
	return b
}

// AddContext adds a context block with mixed elements (text or images).
func (b *BlockBuilder) AddContext(elements ...slack.MixedElement) *BlockBuilder {
	context := NewContextBlock(elements...)
	b.blocks = append(b.blocks, context)
	return b
}

// AddImage adds an image block to the builder.
func (b *BlockBuilder) AddImage(imageURL, altText string) *BlockBuilder {
	image := NewImageBlock(imageURL, altText)
	b.blocks = append(b.blocks, image)
	return b
}

// AddDivider adds a divider block to the builder.
func (b *BlockBuilder) AddDivider() *BlockBuilder {
	divider := NewDividerBlock()
	b.blocks = append(b.blocks, divider)
	return b
}

// Build returns the assembled blocks.
func (b *BlockBuilder) Build() []slack.Block {
	return b.blocks
}
