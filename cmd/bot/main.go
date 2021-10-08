package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/UnderAnder/discord_roll/internal/discord"
	_ "github.com/mattn/go-sqlite3"
)

// Token variable used for command line parameter
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	// We want the seed to change on every launch
	rand.Seed(time.Now().Unix())
}

func main() {
	discordBot, err := discord.NewBot(Token)
	if err != nil {
		log.Fatal(err)
	}
	if err := discordBot.Run(); err != nil {
		log.Fatal(err)
	}
}
