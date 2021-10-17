package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const limit = 10

// topMessage Print leaderboard in response to the text command
func (h *Handler) topMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	top, err := h.top(s)
	if err != nil {
		return
	}

	msg := discordgo.MessageEmbed{
		Description: top,
		Color:       0x006969, // 96 96
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Printf("Failed to response the command %v, %v\n", m.Content, err)
	}
}

// topSlash Print leaderboard in response to the slash command
func (h *Handler) topSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	result, err := h.top(s)
	if err != nil {
		return
	}

	sendRespond(s, i, result)
}

// top Return leaderboard from DB as string
func (h *Handler) top(s *discordgo.Session) (string, error) {
	var sb strings.Builder

	top, err := h.repository.GetTopUsersByScore(limit)
	if err != nil {
		log.Println(err)
		return "", err
	}

	for i, v := range top {
		user, err := s.User(v.DiscordID)
		if err != nil {
			log.Println(err)
			return "", err
		}
		sb.WriteString(fmt.Sprintf("> %d. %s %d\n", i+1, user.Username, v.Score))
	}

	return sb.String(), err
}
