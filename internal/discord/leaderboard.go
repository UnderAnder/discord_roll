package discord

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func leaderboard(b *Bot, s *discordgo.Session) (string, error) {
	var sb strings.Builder

	top, err := b.repository.GetTopUsersByScore(10)
	if err != nil {
		return "", err
	}
	for i, v := range top {
		user, err := s.User(v.Discord_id)
		if err != nil {
			return "", err
		}
		userName := user.Username

		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString(". ")
		sb.WriteString(userName)
		sb.WriteString(" ")
		sb.WriteString(strconv.Itoa(v.Score))
		sb.WriteString("\n")

	}

	return sb.String(), nil
}
