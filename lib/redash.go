package redashslackbot

import (
	"log"
	"os"

	"github.com/justmiles/redash-client"
	"github.com/shomali11/slacker"
)

var client *redash.Client

func init() {
	var err error

	// initialize redash
	client, err = redash.NewClient(os.Getenv("REDASH_URL"), os.Getenv("REDASH_API_KEY"))
	if err != nil {
		log.Fatalf("unable to intialize redash! %s", err)
	}

	// Redash options
	client.DebugEnabled = true

}

// Commands returns available commands
func Commands() map[string]*slacker.CommandDefinition {
	commands := make(map[string]*slacker.CommandDefinition)

	commands["search <search>"] = &slacker.CommandDefinition{
		Description: "Search Redash",
		Example:     "search AWS Billing",
		Handler:     Search,
	}

	commands["show <query>"] = &slacker.CommandDefinition{
		Description: "Show a query chart",
		Example:     "show AWS Billing",
		Handler:     Show,
	}

	return commands
}
