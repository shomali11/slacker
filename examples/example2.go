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

	bot.Default(func(request *slacker.Request, response *slacker.Response) {
		response.Reply("Say what?")
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
