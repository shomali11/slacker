# slacker [![Build Status](https://travis-ci.com/shomali11/slacker.svg?branch=master)](https://travis-ci.com/shomali11/slacker) [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/slacker)](https://goreportcard.com/report/github.com/shomali11/slacker) [![GoDoc](https://godoc.org/github.com/shomali11/slacker?status.svg)](https://godoc.org/github.com/shomali11/slacker) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) 

Built on top of the Slack API [github.com/slack-go/slack](https://github.com/slack-go/slack) with the idea to simplify the Real-Time Messaging feature to easily create Slack Bots, assign commands to them and extract parameters.

## Features

- Easy definitions of commands and their input
- Available bot initialization, errors and default handlers
- Simple parsing of String, Integer, Float and Boolean parameters
- Contains support for `context.Context`
- Built-in `help` command
- Replies can be new messages or in threads
- Supports authorization
- Bot responds to mentions and direct messages
- Handlers run concurrently via goroutines
- Produces events for executed commands
- Full access to the Slack API [github.com/slack-go/slack](https://github.com/slack-go/slack)

## Dependencies

- `commander` [github.com/shomali11/commander](https://github.com/shomali11/commander)
- `slack` [github.com/slack-go/slack](https://github.com/slack-go/slack)

# Install

```
go get github.com/shomali11/slacker
```

# Examples

## Example 1

Defining a command using slacker

```go
package main

import (
	"context"
	"log"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong")
		},
	}

	bot.Command("ping", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 2

Defining a command with an optional description and example. The handler replies to a thread.

```go
package main

import (
	"context"
	"log"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Ping!",
		Example:     "ping",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong", slacker.WithThreadReply(true))
		},
	}

	bot.Command("ping", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 3

Defining a command with a parameter

```go
package main

import (
	"context"
	"log"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Echo a word!",
		Example:     "echo hello",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			word := request.Param("word")
			response.Reply(word)
		},
	}

	bot.Command("echo <word>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 4

Defining a command with two parameters. Parsing one as a string and the other as an integer.
_(The second parameter is the default value in case no parameter was passed or could not parse the value)_

```go
package main

import (
	"context"
	"log"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Repeat a word a number of times!",
		Example:     "repeat hello 10",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			word := request.StringParam("word", "Hello!")
			number := request.IntegerParam("number", 1)
			for i := 0; i < number; i++ {
				response.Reply(word)
			}
		},
	}

	bot.Command("repeat <word> <number>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 5

Defines two commands that display sending errors to the Slack channel. One that replies as a new message. The other replies to the thread.

```go
package main

import (
	"context"
	"errors"
	"log"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	messageReplyDefinition := &slacker.CommandDefinition{
		Description: "Tests errors in new messages",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.ReportError(errors.New("Oops!"))
		},
	}

	threadReplyDefinition := &slacker.CommandDefinition{
		Description: "Tests errors in threads",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.ReportError(errors.New("Oops!"), slacker.WithThreadError(true))
		},
	}

	bot.Command("message", messageReplyDefinition)
	bot.Command("thread", threadReplyDefinition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 6

Send a "Typing" indicator

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Server time!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Typing()

			time.Sleep(time.Second)

			response.Reply(time.Now().Format(time.RFC1123))
		},
	}

	bot.Command("time", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 7

Showcasing the ability to access the [github.com/slack-go/slack](https://github.com/slack-go/slack) API and the Real-Time Messaging Protocol.
_In this example, we are sending a message using RTM and uploading a file using the Slack API._

```go
package main

import (
	"context"
	"log"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Upload a word!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			word := request.Param("word")

			channel := botCtx.Event().Channel
			rtm := botCtx.RTM()
			client := botCtx.Client()

			rtm.SendMessage(rtm.NewOutgoingMessage("Uploading file ...", channel))
			client.UploadFile(slack.FileUploadParameters{Content: word, Channels: []string{channel}})
		},
	}

	bot.Command("upload <word>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 8

Showcasing the ability to leverage `context.Context` to add a timeout

```go
package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Process!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			timedContext, cancel := context.WithTimeout(botCtx.Context(), time.Second)
			defer cancel()

			select {
			case <-timedContext.Done():
				response.ReportError(errors.New("timed out"))
			case <-time.After(time.Minute):
				response.Reply("Processing done!")
			}
		},
	}

	bot.Command("process", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 9

Showcasing the ability to add attachments to a `Reply`

```go
package main

import (
	"context"
	"log"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Echo a word!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			word := request.Param("word")

			attachments := []slack.Attachment{}
			attachments = append(attachments, slack.Attachment{
				Color:      "red",
				AuthorName: "Raed Shomali",
				Title:      "Attachment Title",
				Text:       "Attachment Text",
			})

			response.Reply(word, slacker.WithAttachments(attachments))
		},
	}

	bot.Command("echo <word>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 10

Showcasing the ability to add blocks to a `Reply`

```go
package main

import (
	"context"
	"log"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Echo a word!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			word := request.Param("word")

			attachments := []slack.Block{}
			attachments = append(attachments, slack.NewContextBlock("1",
				slack.NewTextBlockObject("mrkdwn", "Hi!", false, false)),
			)

			response.Reply(word, slacker.WithBlocks(attachments))
		},
	}

	bot.Command("echo <word>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 11

Showcasing the ability to create custom responses via `CustomResponse`

```go
package main

import (
	"log"

	"context"
	"errors"
	"fmt"

	"github.com/shomali11/slacker"
)

const (
	errorFormat = "> Custom Error: _%s_"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.CustomResponse(NewCustomResponseWriter)

	definition := &slacker.CommandDefinition{
		Description: "Custom!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("custom")
			response.ReportError(errors.New("oops"))
		},
	}

	bot.Command("custom", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

// NewCustomResponseWriter creates a new ResponseWriter structure
func NewCustomResponseWriter(botCtx slacker.BotContext) slacker.ResponseWriter {
	return &MyCustomResponseWriter{botCtx: botCtx}
}

// MyCustomResponseWriter a custom response writer
type MyCustomResponseWriter struct {
	botCtx slacker.BotContext
}

// ReportError sends back a formatted error message to the channel where we received the event from
func (r *MyCustomResponseWriter) ReportError(err error, options ...slacker.ReportErrorOption) {
	rtm := r.botCtx.RTM()
	event := r.botCtx.Event()
	rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), event.Channel))
}

// Typing send a typing indicator
func (r *MyCustomResponseWriter) Typing() {
	rtm := r.botCtx.RTM()
	event := r.botCtx.Event()
	rtm.SendMessage(rtm.NewTypingMessage(event.Channel))
}

// Reply send a attachments to the current channel with a message
func (r *MyCustomResponseWriter) Reply(message string, options ...slacker.ReplyOption) {
	rtm := r.botCtx.RTM()
	event := r.botCtx.Event()
	rtm.SendMessage(rtm.NewOutgoingMessage(message, event.Channel))
}
```

## Example 12

Showcasing the ability to toggle the slack Debug option via `WithDebug`

```go
package main

import (
	"context"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>", slacker.WithDebug(true))

	definition := &slacker.CommandDefinition{
		Description: "Ping!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong")
		},
	}

	bot.Command("ping", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 13

Defining a command that can only be executed by authorized users

```go
package main

import (
	"context"
	"log"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	authorizedUsers := []string{"<USER ID>"}

	authorizedDefinition := &slacker.CommandDefinition{
		Description: "Very secret stuff",
		AuthorizationFunc: func(botCtx slacker.BotContext, request slacker.Request) bool {
			return contains(authorizedUsers, botCtx.Event().User)
		},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("You are authorized!")
		},
	}

	bot.Command("secret", authorizedDefinition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func contains(list []string, element string) bool {
	for _, value := range list {
		if value == element {
			return true
		}
	}
	return false
}
```

## Example 14

Adding handlers to when the bot is connected, encounters an error and a fallback for when none of the commands match

```go
package main

import (
	"log"

	"context"
	"fmt"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Init(func() {
		log.Println("Connected!")
	})

	bot.Err(func(err string) {
		log.Println(err)
	})

	bot.DefaultCommand(func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Say what?")
	})

	bot.FallbackEvent(func(event interface{}) {
		fmt.Println(event)
	})

	definition := &slacker.CommandDefinition{
		Description: "help!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("Your own help function...")
		},
	}

	bot.Help(definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 15

Listening to the Commands Events being produced

```go
package main

import (
	"fmt"
	"log"

	"context"

	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Message)
		fmt.Println()
	}
}

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	go printCommandEvents(bot.CommandEvents())

	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong")
		},
	})

	bot.Command("echo <word>", &slacker.CommandDefinition{
		Description: "Echo a word!",
		Example:     "echo hello",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			word := request.Param("word")
			response.Reply(word)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```
