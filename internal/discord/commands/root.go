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

var (
	helpCommand   = map[string]bool{"help": true, "рудз": true, "помощь": true}
	rollCommand   = map[string]bool{"roll": true, "кщдд": true, "ролл": true}
	cityCommand   = map[string]bool{"city": true, "сшен": true, "город": true, "г": true}
	betCommand    = map[string]bool{"bet": true, "иуе": true, "бет": true, "ставка": true}
	topCommand    = map[string]bool{"top": true, "ещз": true, "топ": true, "leaderboard": true, "лидеры": true}
	bottleCommand = map[string]bool{"bottle": true, "ищееду": true, "бутылочка": true}
	scoreCommand  = map[string]bool{"score": true, "ысщку": true, "очки": true}
)

type handler func(*discordgo.Session, *discordgo.MessageCreate)

type Handler struct {
	repository repository.Repository
	eventChan  chan string
}

func NewHandler(r repository.Repository, e chan string) *Handler {
	return &Handler{repository: r, eventChan: e}
}

func (h *Handler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	var handle handler

	switch {
	case helpCommand[command]:
		handle = h.help
	case rollCommand[command]:
		handle = h.roll
	case cityCommand[command]:
		handle = h.city
	case betCommand[command]:
		handle = h.bet
	case topCommand[command]:
		handle = h.top
	case bottleCommand[command]:
		handle = h.bottle
	case scoreCommand[command]:
		handle = h.score
	default:
		log.Printf("Received unrecognized command: %s\n", commandText)
		handle = h.help
	}

	log.Printf("Received command: %s\n", commandText)
	handle(s, m)
}

func getMessageAuthorNick(m *discordgo.MessageCreate) string {
	if m.Member.Nick != "" {
		return m.Member.Nick
	}
	return m.Author.Username
}

func getMemberNick(m *discordgo.Member) string {
	if m.Nick != "" {
		return m.Nick
	}
	return m.User.Username
}
