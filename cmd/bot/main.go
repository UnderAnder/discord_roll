package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/UnderAnder/discord_roll/internal/discord"
)

func main() {
	// flags
	token := flag.String("t", "", "Bot Token")
	dbPath := flag.String("db", "", "Path to DB")
	regCommands := flag.Bool("regcommands", false, "Create Discord slash commands")
	delCommands := flag.Bool("delcommands", false, "Remove Discord slash commands")
	flag.Parse()

	// seed
	rand.Seed(time.Now().Unix())

	// bot instance
	discordBot, err := discord.NewBot(*token, *dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// run
	if err := discordBot.Run(*regCommands, *delCommands); err != nil {
		log.Fatal(err)
	}
}
