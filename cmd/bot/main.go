package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/UnderAnder/discord_roll/internal/config"
	"github.com/UnderAnder/discord_roll/internal/discord"
)

func main() {
	// seed
	rand.Seed(time.Now().Unix())

	// config
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Error getting config: ", err)
	}

	// bot instance
	discordBot, err := discord.NewBot(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// run
	if err := discordBot.Run(); err != nil {
		log.Fatal(err)
	}
}
