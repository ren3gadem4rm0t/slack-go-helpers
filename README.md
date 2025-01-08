# slack-go-helpers

`slack-go-helpers` is a collection of abstractions and helpers designed to simplify building Slack messages using the [`slack-go/slack`](https://github.com/slack-go/slack) library. It includes builders for creating Slack Block Kit messages, attachments, and more.

## Table of Contents

- [Installation](#installation)
- [Packages](#packages)
  - [blockbuilder](#blockbuilder)
    - [BlockBuilder](#blockbuilder)
    - [AttachmentBuilder](#attachmentbuilder)
    - [Color](#color)
  - [deduper](#deduper)
    - [Dedupe](#dedupe)
    - [EvictionPolicy](#evictionpolicy)
    - [Helper Functions](#helper-functions)
- [Examples](#examples)
  - [Basic BlockBuilder Example](#basic-blockbuilder-example)
  - [AttachmentBuilder Example](#attachmentbuilder-example)
  - [Basic Dedupe Example](#basic-dedupe-example)
  - [Custom Eviction Policy Example](#custom-eviction-policy-example)
  - [Automatic Eviction Example](#automatic-eviction-example)
  - [Comprehensive Dedupe and Socket Mode Example](#comprehensive-dedupe-and-socket-mode-example)
  - [Socket Mode Integration Example](#socket-mode-integration-example)
- [License](#license)

## Installation

To install the library, use:

```bash
go get github.com/ren3gadem4rm0t/slack-go-helpers
```

Then, import the packages as needed:

```go
import (
    blockbuilder "github.com/ren3gadem4rm0t/slack-go-helpers/blockbuilder"
    deduper "github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
    aws_helpers "github.com/ren3gadem4rm0t/slack-go-helpers/aws_helpers"
)
```

## Packages

### aws

The `aws` package provides helper functions for extracting useful information from AWS Key IDs, such as determining the AWS Account ID or the resource type based on the key ID prefix.

#### AWSAccountFromAWSKeyID

The `AWSAccountFromAWSKeyID` function decodes an AWS Key ID to extract the associated AWS Account ID.

#### API

- **`AWSAccountFromAWSKeyID(awsKeyID string) (string, error)`**: Decodes the provided AWS Key ID and extracts the 12-digit AWS Account ID. If the decoding fails, it returns an error.

#### Example Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/ren3gadem4rm0t/slack-go-helpers/aws_helpers"
)

func main() {
	keyID := "AKIAEXAMPLE1234567890"
	accountID, err := aws_helpers.AWSAccountFromAWSKeyID(keyID)
	if err != nil {
		log.Fatalf("Error extracting account ID: %v", err)
	}
	fmt.Printf("AWS Account ID: %s\n", accountID)
}
```

#### AWSResourceTypeFromPrefix

The AWSResourceTypeFromPrefix function determines the AWS resource type based on the prefix of an AWS Key ID.

#### API
- **`AWSResourceTypeFromPrefix(awsKeyID string) (string, error)`**: Identifies the resource type based on the first four characters (prefix) of the AWS Key ID. If the prefix is not recognized, it returns “Unknown resource type.

#### Example Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/ren3gadem4rm0t/slack-go-helpers/aws_helpers"
)

func main() {
	keyID := "AKIAEXAMPLE1234567890"
	resourceType, err := aws_helpers.AWSResourceTypeFromPrefix(keyID)
	if err != nil {
		log.Fatalf("Error determining resource type: %v", err)
	}
	fmt.Printf("Resource Type: %s\n", resourceType)
}
```

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

The `Color` package provides pre-defined constants for common Slack message attachment colors as well as an extensive list of additional colors for various contexts.

#### API

- **Primary Colors**:
  - `ColorGood`: Green color (`#36a64f`), typically used for success messages.
  - `ColorWarning`: Orange color (`#ffae42`), typically used for warning messages.
  - `ColorDanger`: Red color (`#ff0000`), typically used for error messages.
  - `ColorInfo`: Blue color (`#439fe0`), typically used for informational messages.

- **Additional Colors**:
  - `ColorRed`: Same as `ColorDanger` (`#ff0000`).
  - `ColorGreen`: Same as `ColorGood` (`#36a64f`).
  - `ColorOrange`: Same as `ColorWarning` (`#ffae42`).
  - `ColorBlue`: Same as `ColorInfo` (`#439fe0`).
  - Other descriptive colors include:
    - `ColorYellow`: Bright yellow (`#ffff00`), ideal for attention-grabbing messages.
    - `ColorPurple`: Purple (`#800080`), for emphasis or special contexts.
    - `ColorPink`: Pink (`#ff69b4`), for playful or informal messages.
    - `ColorGray`: Gray (`#808080`), for neutral or subtle messages.
    - `ColorTeal`: Teal (`#008080`), for calm tones.
    - Full list available in the [blockbuilder/color.go](./blockbuilder/color.go).

#### Example Usage

```go
package main

import (
    "fmt"
    "github.com/ren3gadem4rm0t/slack-go-helpers/blockbuilder"
)

func main() {
    // Example with primary colors
    successAttachment := blockbuilder.NewAttachmentBuilder(blockbuilder.ColorGood).
        AddSection("Operation completed successfully!", true).
        Build()

    // Example with additional colors
    playfulAttachment := blockbuilder.NewAttachmentBuilder(blockbuilder.ColorPink).
        AddSection("Have a wonderful day!", true).
        Build()

    fmt.Println(successAttachment) // JSON for Slack API
    fmt.Println(playfulAttachment) // JSON for Slack API
}
```

### deduper

The `deduper` package provides mechanisms to handle event deduplication and caching. It ensures that duplicate events are ignored, which is particularly useful when dealing with events that might be retried or sent multiple times by Slack.

#### Dedupe

The `Dedupe` struct is responsible for managing event deduplication using an internal cache with configurable eviction policies. It offers both manual and automatic eviction capabilities to maintain cache integrity.

#### API

- **`NewDedupe(sizeLimit int, timeLimit time.Duration, countLimit int, opts ...Option) *Dedupe`**: Initializes a new deduplication handler with specified size, time, and count limits for the cache. Accepts functional options to configure additional behaviors, such as enabling automatic eviction.
  
- **`NewDedupeWithEvictPolicy(evictPolicy *EvictionPolicy, opts ...Option) *Dedupe`**: Initializes a new deduplication handler with a custom eviction policy. Accepts functional options similar to `NewDedupe`.
  
- **`AddEvent(eventID string) bool`**: Checks if an event is a duplicate. If not, it adds the event to the cache and returns `true`. Returns `false` if the event is a duplicate.
  
- **`Middleware(next func(evt *socketmode.Event, client *socketmode.Client)) func(evt *socketmode.Event, client *socketmode.Client)`**: Wraps a Socket Mode event handler with deduplication logic.
  
- **`Size() int`**: Returns the current size of the deduplication cache.
  
- **`Items() map[string]time.Time`**: Returns a copy of the current items in the cache.
  
- **`TriggerEviction()`**: Manually triggers the eviction policy to remove stale or excess events from the cache.
  
- **`StopAutoEviction()`**: Stops the automatic eviction goroutine if it was enabled during initialization.

#### Options

- **`OptionAutoEvict(interval time.Duration) Option`**: Enables automatic eviction with the specified interval. When enabled, a background goroutine periodically applies the eviction policy based on the provided interval.

#### Example Usage

##### **Basic Dedupe Example**

```go
package main

import (
    "fmt"
    "time"

    "github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
)

func main() {
    // Initialize the deduplication handler with size limit, time limit, and count limit
    dedupeHandler := deduper.NewDedupe(1000, 5*time.Minute, 500)

    eventIDs := []string{"event1", "event2", "event3", "event1", "event2"}

    for _, id := range eventIDs {
        if dedupeHandler.AddEvent(id) {
            fmt.Printf("Processed event: %s\n", id)
        } else {
            fmt.Printf("Duplicate event ignored: %s\n", id)
        }
    }

    fmt.Printf("Current cache size: %d\n", dedupeHandler.Size())
}
```

##### **Custom Eviction Policy Example**

```go
package main

import (
    "fmt"
    "time"

    "github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
)

// CustomEvictionPolicy defines a custom eviction policy
type CustomEvictionPolicy struct {
    *deduper.EvictionPolicy
}

// Apply overrides the default Apply method to implement custom logic
func (e *CustomEvictionPolicy) Apply(c *deduper.Cache) {
    // Custom rule: Evict events that start with "temp"
    for id := range c.Items() {
        if len(id) >= 4 && id[:4] == "temp" {
            c.Evict(id)
            fmt.Printf("Evicted event based on custom rule: %s\n", id)
        }
    }

    // Call the original eviction logic
    e.EvictionPolicy.Apply(c)
}

func main() {
    // Initialize the custom eviction policy
    baseEvictionPolicy := deduper.NewEvictionPolicy(1000, 10*time.Minute, 1000)
    customEvictionPolicy := &CustomEvictionPolicy{baseEvictionPolicy}

    // Initialize the deduplication handler with the custom eviction policy
    dedupeHandler := deduper.NewDedupeWithEvictPolicy(customEvictionPolicy.EvictionPolicy)

    // Add events
    events := []string{"temp_event1", "event2", "temp_event3", "event4"}

    for _, id := range events {
        dedupeHandler.AddEvent(id)
    }

    fmt.Printf("Cache size before eviction: %d\n", dedupeHandler.Size())

    // Manually trigger eviction using the ApplyEviction method
    dedupeHandler.TriggerEviction()

    fmt.Printf("Cache size after eviction: %d\n", dedupeHandler.Size())

    // Check remaining events
    for id := range dedupeHandler.Items() {
        fmt.Printf("Remaining event: %s\n", id)
    }
}
```

##### **Automatic Eviction Example**

```go
package main

import (
    "fmt"
    "sync"
    "time"

    "github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
)

func main() {
    // Initialize the deduplication handler with size limit, time limit, and count limit
    // Enable automatic eviction with an interval of 10 seconds
    dedupeHandler := deduper.NewDedupe(1000, 20*time.Second, 1000, deduper.OptionAutoEvict(10*time.Second))
    defer dedupeHandler.StopAutoEviction() // Ensure the eviction goroutine is stopped when main exits

    var wg sync.WaitGroup

    // Include some "temp_*" events to test manual eviction
    eventIDs := []string{
        "event1",
        "event2",
        "temp_event3", // Should be evicted manually
        "event1",       // Duplicate
        "event2",       // Duplicate
        "event4",
        "temp_event5", // Should be evicted manually
        "event6",
    }

    for _, id := range eventIDs {
        wg.Add(1)
        go func(eventID string) {
            defer wg.Done()
            if eventID == "event6" {
                // Wait for 18 seconds before processing event6 to allow some events to expire
                time.Sleep(18 * time.Second)
            }
            if dedupeHandler.AddEvent(eventID) {
                fmt.Printf("Processed event: %s\n", eventID)
            } else {
                fmt.Printf("Duplicate event ignored: %s\n", eventID)
            }
        }(id)
    }

    wg.Wait()

    fmt.Printf("Final cache size: %d\n", dedupeHandler.Size())

    // Manually trigger eviction to remove "temp_*" events
    fmt.Println("Manually triggering eviction...")
    dedupeHandler.TriggerEviction()

    // After manual eviction, check cache size and contents
    fmt.Printf("Cache size after manual eviction: %d\n", dedupeHandler.Size())
    fmt.Println("Remaining events after manual eviction:")
    for id := range dedupeHandler.Items() {
        fmt.Printf("- %s\n", id)
    }

    // Wait for automatic eviction to occur
    fmt.Println("Waiting for automatic eviction to occur...")
    time.Sleep(15 * time.Second)

    // Final cache size after automatic eviction
    fmt.Printf("Cache size after automatic eviction: %d\n", dedupeHandler.Size())
    fmt.Println("Remaining events after automatic eviction:")
    for id := range dedupeHandler.Items() {
        fmt.Printf("- %s\n", id)
    }
}
```

#### EvictionPolicy

The `EvictionPolicy` struct defines the rules for removing stale entries from the cache. It supports eviction based on time limits and count limits.

#### API

- **`NewEvictionPolicy(sizeLimit int, timeLimit time.Duration, countLimit int) *EvictionPolicy`**: Creates a new eviction policy with specified size, time, and count limits.
  
- **`Apply(c *Cache)`**: Applies the eviction policy to the given cache, removing items that exceed the defined limits.

#### Example Usage

```go
package main

import (
    "time"

    "github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
)

func main() {
    // Create an eviction policy with a size limit of 1000, time limit of 10 minutes, and count limit of 500
    evictionPolicy := deduper.NewEvictionPolicy(1000, 10*time.Minute, 500)

    // Initialize the cache with the eviction policy
    cache := deduper.NewCache(evictionPolicy)

    // Add events to the cache
    cache.Add("event_1")
    cache.Add("event_2")

    // Check if an event exists
    if cache.Has("event_1") {
        println("Event 1 is in the cache.")
    }

    // Get the current size of the cache
    println("Cache size:", cache.Size())
}
```

#### Helper Functions

The `deduper` package also includes helper functions to extract stable event IDs from Slack Socket Mode events, ensuring accurate deduplication.

#### API

- **`ExtractEventIDFromSocketMode(evt *socketmode.Event) (string, error)`**: Extracts a stable event ID from a Socket Mode event for deduplication purposes. It prefers `client_msg_id` for message events and falls back to hashing key event fields for other event types.
  
- **`ExtractEnvelopeIDFromSocketMode(evt *socketmode.Event) (string, error)`**: Extracts the `envelope_id` from a Socket Mode event, which can also be used for deduplication.

#### Example Usage

```go
package main

import (
    "fmt"

    "github.com/ren3gadem4rm0t/slack-go-helpers/deduper"
    "github.com/slack-go/slack/socketmode"
)

func main() {
    // Assume you have a socketmode.Event named evt
    var evt *socketmode.Event

    // Extract the event ID
    eventID, err := deduper.ExtractEventIDFromSocketMode(evt)
    if err != nil {
        fmt.Printf("Error extracting event ID: %v\n", err)
        return
    }

    fmt.Printf("Extracted Event ID: %s\n", eventID)

    // Alternatively, extract the envelope ID
    envelopeID, err := deduper.ExtractEnvelopeIDFromSocketMode(evt)
    if err != nil {
        fmt.Printf("Error extracting envelope ID: %v\n", err)
        return
    }

    fmt.Printf("Extracted Envelope ID: %s\n", envelopeID)
}
```

---

## Examples

### Basic BlockBuilder Example

See [`examples/blockbuilder/basic/main.go`](./examples/blockbuilder/basic/main.go) for a simple usage of the `BlockBuilder`.

### AttachmentBuilder Example

See [`examples/blockbuilder/attachment/main.go`](./examples/blockbuilder/attachment/main.go) for an example of creating attachments with blocks.

### Basic Dedupe Example

See [`examples/deduper/basic/main.go`](./examples/deduper/basic/main.go) for a simple usage of the `Dedupe` struct.

### Custom Eviction Policy Example

See [`examples/deduper/custom_eviction/main.go`](./examples/deduper/custom_eviction/main.go) for an example of implementing a custom eviction policy with `Dedupe`.

### Automatic Eviction Example

See [`examples/deduper/auto_eviction/main.go`](./examples/deduper/auto_eviction/main.go) for an example of enabling and using automatic eviction in the `Dedupe` struct.

### Comprehensive Dedupe and Socket Mode Example

See [`examples/deduper/comprehensive/main.go`](./examples/deduper/comprehensive/main.go) for a comprehensive example integrating `Dedupe` with Slack Socket Mode.

### Socket Mode Integration Example

See [`examples/deduper/socketmode/main.go`](./examples/deduper/socketmode/main.go) for an example of integrating `Dedupe` with Slack Socket Mode.

### AWSAccountFromAWSKeyID Example

See [`examples/aws_helpers/account_from_keyid/main.go`](./examples/aws_helpers/account_from_keyid/main.go) for an example of using AWSAccountFromAWSKeyID.

### AWSResourceTypeFromPrefix Example

See [`examples/aws_helpers/resource_type_from_prefix/main.go`](./examples/aws_helpers/resource_type_from_prefix/main.go) for an example of using AWSResourceTypeFromPrefix.

---

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
