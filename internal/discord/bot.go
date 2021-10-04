package discord

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/UnderAnder/discord_roll/internal/repository"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	bot        *discordgo.Session
	repository repository.Repository
}

// constructor
func NewBot(bot *discordgo.Session, repository repository.Repository) *Bot {
	return &Bot{bot: bot, repository: repository}
}

func (b *Bot) Start() error {
	// Register the messageCreate func as a callback for MessageCreate events.
	b.bot.AddHandler(b.messageCreate)

	// In this example, we only care about receiving message events.
	b.bot.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err := b.bot.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
		return err
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	b.bot.Close()

	return nil
}
