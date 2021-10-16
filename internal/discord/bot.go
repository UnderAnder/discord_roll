package discord

import (
	"errors"
	"github.com/UnderAnder/discord_roll/internal/discord/commands"
	"github.com/UnderAnder/discord_roll/internal/discord/reactions"
	"github.com/UnderAnder/discord_roll/internal/repository"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Bot listens to Discord and performs the various actions
type Bot struct {
	discord    *discordgo.Session
	repository repository.Repository
}

var GuildID = "" // Register slash commands globally

// NewBot configures a Bot and returns it.
func NewBot(token, dbPath string) (*Bot, error) {
	db, err := repository.New(dbPath)
	if err != nil {
		log.Fatal("error creating DB session,", err)
		return nil, err
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Panic("error creating Discord session,", err)
		return nil, err
	}

	// Create Event chans
	events := make(chan string)

	commandsHandler := commands.NewHandler(db, events)
	// Register text commands handler
	dg.AddHandler(commandsHandler.HandleMessage)
	// Register slash commands handler
	dg.AddHandler(commandsHandler.HandleInteraction)

	// Register reactionAdd handler
	reactionsHandler := reactions.NewHandler(db, events)
	dg.AddHandler(reactionsHandler.HandleAdd)

	return &Bot{
		discord:    dg,
		repository: db,
	}, nil
}

// Run starts the bot, listens for a halt signal, and shuts down when the halt is received.
func (b *Bot) Run(regCommands, delCommands bool) error {
	if err := b.Start(regCommands); err != nil {
		return errors.New("failed to start bot")
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Received stop signal, shutting down...")
	b.Stop(delCommands)
	return nil
}

// Start opens the connection to the discord web socket.
func (b *Bot) Start(regCommands bool) error {
	log.Println("Starting bot...")
	if err := b.repository.Open(); err != nil {
		return errors.New("failed to open repository")
	}

	log.Println("Opening connection to Discord...")
	if err := b.discord.Open(); err != nil {
		return errors.New("failed to open web socket connection to Discord")
	}
	log.Println("Connection to Discord established.")

	// create discord slash commands
	if regCommands {
		b.createCommands()
	}
	return nil
}

// Stop gracefully shuts down the bot.
func (b *Bot) Stop(delCommands bool) {
	log.Println("Stopping bot...")

	// removing slash commands before exit
	if delCommands {
		b.delCommands()
	}

	log.Println("Closing connection to Discord...")
	err := b.discord.Close()
	if err != nil {
		log.Println(err)
	}

	log.Println("Closing connection to DB...")
	err = b.repository.Close()
	if err != nil {
		log.Printf("Error closing store session: %v\n", err)
	}
}

// delCommands removes slash commands from discord
func (b *Bot) delCommands() {
	log.Println("Removing slash commands...")
	registeredCommands, _ := b.discord.ApplicationCommands(b.discord.State.User.ID, GuildID)
	for _, v := range registeredCommands {
		err := b.discord.ApplicationCommandDelete(b.discord.State.User.ID, GuildID, v.ID)
		log.Printf("remove command: %v id: %v\n", v.Name, v.ID)
		if err != nil {
			log.Printf("Cannot delete '%v' command: %v\n", v.Name, err)
		}
	}
}

// createCommands creates slash commands in discord
func (b *Bot) createCommands() {
	for _, v := range commands.SlashCommands {
		_, err := b.discord.ApplicationCommandCreate(b.discord.State.User.ID, GuildID, v)
		log.Printf("Create command %v\n", v.Name)
		if err != nil {
			log.Printf("Cannot create '%v' command: %v", v.Name, err)
		}
	}
	log.Println("Slash commands will be available on the Discord server in a few minutes")
}
