package main

import (
	"context"
	"fmt"
	"log"
	"os"

	lib "github.com/justmiles/redash-slackbot/lib"
	"github.com/shomali11/slacker"
)

func main() {

	var enableDebugging bool
	if os.Getenv("DEBUG") != "" {
		enableDebugging = true
	}

	bot := slacker.NewClient(os.Getenv("SLACK_TOKEN"), slacker.WithDebug(enableDebugging))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for usage, definition := range lib.Commands() {
		bot.Command(usage, definition)
	}

	fmt.Println("Starting Redash Slackbot")
	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
