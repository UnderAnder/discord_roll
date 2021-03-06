package commands

import (
	"log"
	"strings"

	"github.com/UnderAnder/discord_roll/internal/config"
	"github.com/UnderAnder/discord_roll/internal/repository"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	CommandKeyword = "!"
)

// SlashCommands list of slash commands for registration on Discord
var SlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "top",
		Description: "Leaderboard",
	},
	{
		Name:        "score",
		Description: "Show your score",
	},
	{
		Name:        "roll",
		Description: "Generate a random number",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "max",
				Description: "Will generate a random number between 1 and 'max'. Default: 100",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "quantity",
				Description: "How many times. Default: 1",
				Required:    false,
			},
		},
	},
	{
		Name:        "duel",
		Description: "make roll against opponent",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "opponent",
				Description: "opponent",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "bet",
				Description: "bet",
				Required:    true,
			},
		},
	},
	{
		Name:        "bet",
		Description: "Make a bet",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "bet",
				Description: "bet",
				Required:    true,
			},
		},
	},
	{
		Name:        "city",
		Description: "guess the city",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "city",
				Description: "city",
				Required:    true,
			},
		},
	},
}

// list of text commands for use on a guild channels
var (
	helpCommand  = map[string]bool{"help": true, "????????": true, "????????????": true}
	rollCommand  = map[string]bool{"roll": true, "????????": true, "????????": true}
	duelCommand  = map[string]bool{"duel": true, "????????": true, "??????????": true}
	cityCommand  = map[string]bool{"city": true, "c": true, "????????": true, "??????????": true, "??": true}
	betCommand   = map[string]bool{"bet": true, "??????": true, "??????": true, "????????????": true}
	topCommand   = map[string]bool{"top": true, "??????": true, "??????": true, "leaderboard": true, "????????????": true}
	scoreCommand = map[string]bool{"score": true, "??????????": true, "????????": true}
)

type handlerMessage func(*discordgo.Session, *discordgo.MessageCreate)
type handlerInteraction func(*discordgo.Session, *discordgo.InteractionCreate)

type Handler struct {
	repository repository.Repository
	eventChan  chan string
	cfg        *config.Config
	localizer  *i18n.Localizer
}

func NewHandler(r repository.Repository, e chan string, cfg *config.Config, loc *i18n.Localizer) *Handler {
	return &Handler{repository: r, eventChan: e, cfg: cfg, localizer: loc}
}

// HandleMessage text command handler
func (h *Handler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages from bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	text := strings.TrimSpace(strings.ToLower(m.Content))

	if !strings.HasPrefix(text, CommandKeyword) {
		return
	}

	// check guild is allowed
	if h.cfg.Bot.GuildID != "" && h.cfg.Bot.GuildID != m.GuildID {
		log.Println("Get command from not allowed guild")
		return
	}

	commandText := strings.TrimSpace(text[len(CommandKeyword):])
	if commandText == "" {
		return
	}
	commands := strings.Fields(commandText)
	command := commands[0]

	var handle handlerMessage

	switch {
	case helpCommand[command]:
		handle = h.help
	case rollCommand[command]:
		handle = h.rollMessage
	case duelCommand[command]:
		handle = h.duelMessage
	case cityCommand[command]:
		handle = h.cityMessage
	case betCommand[command]:
		handle = h.betMessage
	case topCommand[command]:
		handle = h.topMessage
	case scoreCommand[command]:
		handle = h.scoreMessage
	default:
		log.Printf("Received unrecognized command: %s\n", commandText)
		handle = h.help
	}

	log.Printf("Received command: %s\n", commandText)
	handle(s, m)
}

// HandleInteraction slash commands handler
func (h *Handler) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// check guild is allowed
	if h.cfg.Bot.GuildID != "" && h.cfg.Bot.GuildID != i.GuildID {
		log.Println("Get interaction from not allowed guild")
		return
	}

	var handle handlerInteraction

	switch i.ApplicationCommandData().Name {
	case "top":
		handle = h.topSlash
	case "score":
		handle = h.scoreSlash
	case "roll":
		handle = h.rollSlash
	case "duel":
		handle = h.duelSlash
	case "bet":
		handle = h.betSlash
	case "city":
		handle = h.citySlash
	default:
		log.Panicln("UNREGISTERED COMMAND")
		return
	}

	log.Printf("Received slash command: %v", i.ApplicationCommandData().Name)
	handle(s, i)
}

// interactionUserID check interaction came from channel or DM and return UserID
func interactionUserID(i *discordgo.InteractionCreate) string {
	var userID string
	switch {
	case i.Member != nil:
		userID = i.Member.User.ID
	case i.User != nil:
		userID = i.User.ID
	default:
		log.Panicln("Can't get userID")
	}
	return userID
}

func sendRespond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) bool {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		log.Printf("Failed to response the command %v, %v\n", i.ApplicationCommandData().Name, err)
		return false
	}
	return true
}

func sendMessageReply(s *discordgo.Session, m *discordgo.MessageCreate, msg string) (*discordgo.Message, bool) {
	msgInstance, err := s.ChannelMessageSendReply(m.ChannelID, msg, m.Message.Reference())
	if err != nil {
		log.Printf("Failed to response the command %v, %v\n", m.Content, err)
		return nil, false
	}
	return msgInstance, true
}
