package main

import (
	"log"

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

	bot.DefaultCommand(func(request *slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Say what?")
	})

	bot.DefaultEvent(func(event interface{}) {
		fmt.Println(event)
	})

	bot.Help(func(request *slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Your own help function...")
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
