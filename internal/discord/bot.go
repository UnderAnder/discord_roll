package discord

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/UnderAnder/discord_roll/internal/discord/commands"
	"github.com/UnderAnder/discord_roll/internal/repository"
	"github.com/bwmarrin/discordgo"
)

// Bot listens to Discord and performs the various actions
type Bot struct {
	discord    *discordgo.Session
	repository repository.Repository
}

// NewBot configures a Bot and returns it.
func NewBot(token string) (*Bot, error) {
	// NOTICE hardcoded path
	db, err := repository.New("./db.sqlite3")
	if err != nil {
		return nil, err
	}
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return nil, err
	}

	commandsHandler := commands.NewHandler(db)

	dg.AddHandler(commandsHandler.Handle)

	return &Bot{
		discord:    dg,
		repository: db,
	}, nil
}

// Start opens the connection to the discord web socket.
func (b *Bot) Start() error {
	log.Println("Starting bot...")
	if err := b.repository.Open(); err != nil {
		return errors.New("failed to open repository")
	}

	log.Println("Opening connection to Discord...")
	if err := b.discord.Open(); err != nil {
		return errors.New("failed to open web socket connection to Discord")
	}
	log.Println("Connection to Discord established.")
	return nil
}

// Run starts the bot, listens for a halt signal, and shuts down when the halt is received.
func (b *Bot) Run() error {
	if err := b.Start(); err != nil {
		return errors.New("failed to start bot")
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Received stop signal, shutting down...")
	b.Stop()
	return nil
}

// Stop gracefully shuts down the bot.
func (b *Bot) Stop() {
	log.Println("Stopping bot...")
	log.Println("Closing connection to Discord...")
	err := b.discord.Close()
	if err != nil {
		log.Fatal("failed ", err)
	}

	err = b.repository.Close()
	if err != nil {
		log.Printf("Error closing store session: %v", err)
	}
}
