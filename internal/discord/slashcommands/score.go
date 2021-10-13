package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
)

// Returns user score
func (h *Handler) score(s *discordgo.Session, i *discordgo.InteractionCreate) {
	score, err := h.repository.GetScore(i.User.ID)
	if err != nil {
		log.Println(err)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: i.User.Username + " у тебя " + strconv.Itoa(score) + " очков",
		},
	})
	if err != nil {
		log.Printf("Failed to handle the command %v, %v\n", i.ApplicationCommandData().Name, err)
		return
	}
}
