package commands

import (
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) bottle(s *discordgo.Session, m *discordgo.MessageCreate) {
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err != nil {
		log.Println("error can't get guild members", err)
		return
	}

	randMember := members[rand.Intn(len(members))] //nolint:gosec
	for randMember.User.ID == m.Author.ID {
		randMember = members[rand.Intn(len(members))] //nolint:gosec
	}

	if _, err := s.ChannelMessageSend(m.ChannelID, getMessageAuthorNick(m)+" :kiss: "+getMemberNick(randMember)); err != nil {
		log.Println(err)
	}
}
