package commands

import (
	"log"
	"strings"

	"github.com/UnderAnder/discord_roll/internal/repository"
	"github.com/bwmarrin/discordgo"
)

const (
	CommandKeyword = "!"
)

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
}

var (
	helpCommand  = map[string]bool{"help": true, "рудз": true, "помощь": true}
	rollCommand  = map[string]bool{"roll": true, "кщдд": true, "ролл": true}
	cityCommand  = map[string]bool{"city": true, "сшен": true, "город": true, "г": true}
	betCommand   = map[string]bool{"bet": true, "иуе": true, "бет": true, "ставка": true}
	topCommand   = map[string]bool{"top": true, "ещз": true, "топ": true, "leaderboard": true, "лидеры": true}
	scoreCommand = map[string]bool{"score": true, "ысщку": true, "очки": true}
)

type handlerMessage func(*discordgo.Session, *discordgo.MessageCreate)
type handlerInteraction func(*discordgo.Session, *discordgo.InteractionCreate)

type Handler struct {
	repository repository.Repository
	eventChan  chan string
}

func NewHandler(r repository.Repository, e chan string) *Handler {
	return &Handler{repository: r, eventChan: e}
}

func (h *Handler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	text := strings.TrimSpace(strings.ToLower(m.Content))

	if !strings.HasPrefix(text, CommandKeyword) {
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
	case cityCommand[command]:
		handle = h.city
	case betCommand[command]:
		handle = h.bet
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

func (h *Handler) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Received slash command: %v", i.ApplicationCommandData().Name)

	var handle handlerInteraction

	switch i.ApplicationCommandData().Name {
	case "top":
		handle = h.topSlash
	case "score":
		handle = h.scoreSlash
	case "roll":
		handle = h.rollSlash
	default:
		log.Panicln("UNREGISTRED COMMAND")
		return
	}
	handle(s, i)
}
