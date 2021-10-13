package slashcommands

import (
	"github.com/UnderAnder/discord_roll/internal/repository"
	"github.com/bwmarrin/discordgo"
	"log"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "top",
		Description: "Leaderboard",
	},
	{
		Name:        "score",
		Description: "Show your score",
	},
}

type handler func(*discordgo.Session, *discordgo.InteractionCreate)

type Handler struct {
	repository repository.Repository
	eventChan  chan string
}

func NewHandler(r repository.Repository, e chan string) *Handler {
	return &Handler{repository: r, eventChan: e}
}

func (h *Handler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Received slash command: %v", i.ApplicationCommandData().Name)

	var handle handler

	switch i.ApplicationCommandData().Name {
	case "top":
		handle = h.top
	case "score":
		handle = h.score
	default:
		log.Println("ERROR")
	}
	handle(s, i)
}
