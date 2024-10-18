# slack-go-helpers

`slack-go-helpers` is a collection of abstractions and helpers designed to simplify building Slack messages using the [`slack-go/slack`](https://github.com/slack-go/slack) library. It includes builders for creating Slack Block Kit messages, attachments, and more.

## Table of Contents

- [Installation](#installation)
- [Modules](#modules)
  - [blockbuilder](#blockbuilder)
    - [BlockBuilder](#blockbuilder)
    - [AttachmentBuilder](#attachmentbuilder)
    - [Color](#color)
- [Examples](#examples)
- [License](#license)

## Installation

To install the library, use:

```bash
go get github.com/ren3gadem4rm0t/slack-go-helpers
```

Then, import the packages as needed:

```go
import (
    "github.com/ren3gadem4rm0t/slack-go-helpers/blockbuilder"
)
```

## Modules

### blockbuilder

The `blockbuilder` package provides abstractions for building Slack Block Kit messages and attachments with clean, type-safe APIs. The package also includes helper utilities for managing colors in Slack messages.

#### BlockBuilder

The `BlockBuilder` is responsible for creating blocks for Slack messages. Blocks are the core building components of a Slack message. You can add sections, context, dividers, images, and actions to your message.

#### API

- `NewBlockBuilder() *BlockBuilder`: Creates a new `BlockBuilder`.
- `AddSection(text string, markdown bool) *BlockBuilder`: Adds a section block with either Markdown or plain text.
- `AddActions(elements ...slack.BlockElement) *BlockBuilder`: Adds action elements (e.g., buttons).
- `AddContext(elements ...slack.MixedElement) *BlockBuilder`: Adds a context block with mixed elements (text or images).
- `AddImage(imageURL, altText string) *BlockBuilder`: Adds an image block.
- `AddDivider() *BlockBuilder`: Adds a divider block.
- `Build() []slack.Block`: Returns the assembled blocks.

#### Example Usage

```go
package main

import (
    "fmt"
    "github.com/ren3gadem4rm0t/slack-go-helpers/blockbuilder"
    "github.com/slack-go/slack"
)

func main() {
    builder := blockbuilder.NewBlockBuilder().
        AddSection("Hello from BlockBuilder!", true).
        AddDivider().
        AddActions(
            blockbuilder.NewButton("btn_1", "Click Me", "value_1"),
            blockbuilder.NewButton("btn_2", "Cancel", "value_2"),
        )

    blocks := builder.Build()

    fmt.Println(blocks) // JSON or Slack API usage
}
```

#### AttachmentBuilder

The `AttachmentBuilder` is used for creating Slack message attachments. Attachments allow you to create colored sections with embedded blocks.

#### API

- `NewAttachmentBuilder(color string) *AttachmentBuilder`: Creates a new `AttachmentBuilder` with a color.
- `AddSection(text string, markdown bool) *AttachmentBuilder`: Adds a section block to the attachment.
- `AddActions(elements ...slack.BlockElement) *AttachmentBuilder`: Adds action elements to the attachment.
- `AddContext(elements ...slack.MixedElement) *AttachmentBuilder`: Adds a context block to the attachment.
- `AddImage(imageURL, altText string) *AttachmentBuilder`: Adds an image block to the attachment.
- `AddDivider() *AttachmentBuilder`: Adds a divider block to the attachment.
- `AddBlock(block slack.Block) *AttachmentBuilder`: Adds a custom block to the attachment.
- `AddBlocksFromBuilder(builder *BlockBuilder) *AttachmentBuilder`: Adds blocks from a `BlockBuilder` instance to the attachment.
- `Build() slack.Attachment`: Returns the assembled attachment.

#### Example Usage of `AddBlock`

```go
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
```

#### Example Usage of `AddBlocksFromBuilder`

```go
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
```

#### Color

The `Color` package provides pre-defined constants for common Slack message attachment colors such as success, warning, danger, and informational colors.

#### API

- `ColorGood`: Green color (`#36a64f`), typically used for success messages.
- `ColorWarning`: Orange color (`#ffae42`), typically used for warning messages.
- `ColorDanger`: Red color (`#ff0000`), typically used for error messages.
- `ColorInfo`: Blue color (`#439fe0`), typically used for informational messages.

#### Example Usage

```go
package main

import (
    "github.com/ren3gadem4rm0t/slack-go-helpers/blockbuilder"
)

func main() {
    // Using the color constants to set attachment colors
    attachment := blockbuilder.NewAttachmentBuilder(blockbuilder.ColorGood).
        AddSection("Everything is working as expected!", true).
        Build()

    fmt.Println(attachment) // JSON or Slack API usage
}
```

## Examples

### Basic BlockBuilder Example

See [`examples/blockbuilder/basic/main.go`](./examples/blockbuilder/basic/main.go) for a simple usage of the `BlockBuilder`.

### Attachment Example

See [`examples/blockbuilder/attachment/main.go`](./examples/blockbuilder/attachment/main.go) for an example of creating attachments with blocks.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
