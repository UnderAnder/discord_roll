package discord

import (
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func game_bottle(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.Member, error) {
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err != nil {
		log.Println("error can't get guild members", err)
		return nil, err
	}

	randMember := members[rand.Intn(len(members))]
	for randMember.User.ID == m.Author.ID {
		randMember = members[rand.Intn(len(members))]
	}

	return randMember, nil
}
