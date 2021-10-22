package commands

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) help(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpCommand, _ := h.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "help.msg",
	})

	msg := discordgo.MessageEmbed{
		Description: helpCommand,
		Color:       0x006969, // 96 96
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Println(err)
	}
}
