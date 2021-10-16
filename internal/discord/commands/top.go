package commands

import (
	"log"
	"strconv"
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
	top, err := h.top(s)
	if err != nil {
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: top,
		},
	})
	if err != nil {
		log.Printf("Failed to response the command %v, %v\n", i.ApplicationCommandData().Name, err)
		return
	}
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

		sb.WriteString("> ")
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString(". ")
		sb.WriteString(user.Username)
		sb.WriteString(" ")
		sb.WriteString(strconv.Itoa(v.Score))
		sb.WriteString("\n")
	}

	return sb.String(), err
}
