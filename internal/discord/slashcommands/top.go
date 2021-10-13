package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
)

// Show top10 leaderboard
func (h *Handler) top(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var sb strings.Builder
	//NOTICE: hardcoded
	top, err := h.repository.GetTopUsersByScore(10)
	if err != nil {
		log.Println(err)
		return
	}
	for i, v := range top {
		user, err := s.User(v.DiscordID)
		if err != nil {
			log.Println(err)
			return
		}
		userName := user.Username
		sb.WriteString("> ")
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString(". ")
		sb.WriteString(userName)
		sb.WriteString(" ")
		sb.WriteString(strconv.Itoa(v.Score))
		sb.WriteString("\n")
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
		},
	})
	if err != nil {
		log.Printf("Failed to handle the command %v, %v\n", i.ApplicationCommandData().Name, err)
		return
	}
}
