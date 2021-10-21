package discord

import (
	"errors"
	"github.com/UnderAnder/discord_roll/internal/config"
	"github.com/UnderAnder/discord_roll/internal/discord/commands"
	"github.com/UnderAnder/discord_roll/internal/discord/reactions"
	"github.com/UnderAnder/discord_roll/internal/repository"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Bot listens to Discord and performs the various actions
type Bot struct {
	discord    *discordgo.Session
	repository repository.Repository
	cfg        *config.Config
}

// NewBot configures a Bot and returns it.
func NewBot(cfg *config.Config) (*Bot, error) {
	db, err := repository.New(cfg.Repository.Sqlite.Location)
	if err != nil {
		log.Fatal("error creating DB session,", err)
		return nil, err
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		log.Panic("error creating Discord session,", err)
		return nil, err
	}

	// Create Event chan
	events := make(chan string)

	commandsHandler := commands.NewHandler(db, events, cfg)
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
		cfg:        cfg,
	}, nil
}

// Run starts the bot, listens for a halt signal, and shuts down when the halt is received.
func (b *Bot) Run() error {
	if err := b.Start(); err != nil {
		return err
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Received stop signal, shutting down...")
	b.Stop()
	return nil
}

// Start opens the connection to the discord web socket.
func (b *Bot) Start() error {
	log.Println("Starting bot...")
	if err := b.repository.Open(); err != nil {
		return errors.New("failed to open repository")
	}

	log.Println("Opening connection to Discord...")
	retry(b.discord.Open)
	log.Println("Connection to Discord established.")

	// create discord slash commands
	if b.cfg.Bot.RegCommands {
		b.createCommands()
	}
	return nil
}

// Stop gracefully shuts down the bot.
func (b *Bot) Stop() {
	log.Println("Stopping bot...")

	// removing slash commands before exit
	if b.cfg.Bot.DelCommands {
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

// delCommands removes slash commands from Discord
func (b *Bot) delCommands() {
	guildID := b.cfg.Bot.GuildID
	switch guildID {
	case "":
		log.Println("Removing slash commands globally...")
	default:
		log.Printf("Removing slash commands from guild %v...\n", guildID)
	}
	registeredCommands, _ := b.discord.ApplicationCommands(b.discord.State.User.ID, guildID)
	for _, v := range registeredCommands {
		err := b.discord.ApplicationCommandDelete(b.discord.State.User.ID, guildID, v.ID)
		log.Printf("Remove command: %v\n", v.Name)
		if err != nil {
			log.Printf("Cannot delete '%v' command: %v\n", v.Name, err)
		}
	}
}

// createCommands creates slash commands in Discord
func (b *Bot) createCommands() {
	guildID := b.cfg.Bot.GuildID
	switch guildID {
	case "":
		log.Println("Creating slash commands globally...")
	default:
		log.Printf("Creating slash commands on guild %v...\n", guildID)
	}

	for _, v := range commands.SlashCommands {
		_, err := b.discord.ApplicationCommandCreate(b.discord.State.User.ID, guildID, v)
		log.Printf("Create command %v\n", v.Name)
		if err != nil {
			log.Printf("Cannot create '%v' command: %v", v.Name, err)
		}
	}
	log.Println("Slash commands will be available in a few minutes")
}

// retry to open the Discord connection until it is established
func retry(f func() error) {
	wait := time.Duration(1)
	for {
		err := f()
		if err == nil {
			return
		}
		if errors.Is(err, discordgo.ErrWSAlreadyOpen) {
			log.Println("Websocket already exists, no need to reconnect")
			return
		}
		time.Sleep(wait)
		wait *= 2
		if wait > 600 {
			wait = 600
		}
		log.Println("retrying after error:", err)
	}
}
