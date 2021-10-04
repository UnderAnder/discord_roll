package discord

import (
	"github.com/bwmarrin/discordgo"
)

func getMessageAuthorNick(m *discordgo.MessageCreate) string {
	if m.Member.Nick != "" {
		return m.Member.Nick
	} else {
		return m.Author.Username
	}
}

func getMemberNick(m *discordgo.Member) string {
	if m.Nick != "" {
		return m.Nick
	} else {
		return m.User.Username
	}
}
