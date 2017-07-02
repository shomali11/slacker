package main

import (
	"log"

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

	bot.Default(func(request *slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Say what?")
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
