package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/UnderAnder/discord_roll/internal/discord"
)

func main() {
	const guildIDusage = "ID of Discord guild (server) to register or remove slash commands. " +
		"If not specified, slash commands will be global"
	// flags
	token := flag.String("t", "", "Bot Token")
	dbPath := flag.String("db", "./data/sqlite/bot.sqlite3", "Path to DB")
	regCommands := flag.Bool("regcommands", false, "Create Discord slash commands")
	delCommands := flag.Bool("delcommands", false, "Remove Discord slash commands")
	guildID := flag.String("guild", "", guildIDusage)
	flag.Parse()

	// seed
	rand.Seed(time.Now().Unix())

	// bot instance
	discordBot, err := discord.NewBot(*token, *dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// run
	if err := discordBot.Run(*regCommands, *delCommands, *guildID); err != nil {
		log.Fatal(err)
	}
}
