package commands

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const limit = 10

func (h *Handler) top(s *discordgo.Session, m *discordgo.MessageCreate) {
	var sb strings.Builder

	top, err := h.repository.GetTopUsersByScore(limit)
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

		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString(". ")
		sb.WriteString(userName)
		sb.WriteString(" ")
		sb.WriteString(strconv.Itoa(v.Score))
		sb.WriteString("\n")
	}
	msg := discordgo.MessageEmbed{
		Description: sb.String(),
		Color:       0x006969, // 96 96
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Println(err)
	}
}
