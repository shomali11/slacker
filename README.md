# slacker [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/slacker)](https://goreportcard.com/report/github.com/shomali11/slacker)

Built on top of the Slack API https://github.com/nlopes/slack with the idea to simplify the Real-Time Messaging protocol to easily build Slack Bots.

# Examples

## Example 1

Defining a command using slacker

```
package main

import (
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("ping", "Ping!", func(request *slacker.Request, response *slacker.Response) {
		response.Reply("Pong")
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 2

Adding handlers to when the bot is connected and encounters an error

```
package main

import (
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Init(func() {
		log.Println("Connected!")
	})

	bot.Err(func(err string) {
		log.Println(err)
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 3

Defining a command with a parameter

```
package main

import (
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("echo <word>", "Echo a word!", func(request *slacker.Request, response *slacker.Response) {
		word := request.Param("word")
		response.Reply(word)
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 4

Defining a command with two parameters. Parsing one as a string and the other as an integer. 
_(The second parameter is the default value in case no parameter was passed or could not parse the value)_

```
package main

import (
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("repeat <word> <number>", "Repeat a word a number of times!", func(request *slacker.Request, response *slacker.Response) {
		word := request.StringParam("word", "Hello!")
		number := request.IntegerParam("number", 1)
		for i := 0; i < number; i++ {
			response.Reply(word)
		}
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 5

Showcasing the ability to access the https://github.com/nlopes/slack API. 
_In this example, we upload a file using the Slack API._

```
package main

import (
	"github.com/nlopes/slack"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("upload <word>", "Upload a word!", func(request *slacker.Request, response *slacker.Response) {
		word := request.Param("word")
		channel := request.Event.Channel
		bot.Client.UploadFile(slack.FileUploadParameters{Content: word, Channels: []string{channel}})
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
```
