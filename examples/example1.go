package main

import (
"github.com/shomali11/slacker"
"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("ping", "Ping!", func(request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("pong")
	})

	ctx := context.Background()
	// Call cancel() for graceful shutdown
	ctx, cancel := context.WithCancel(ctx)

	if err := bot.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}