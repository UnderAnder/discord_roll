package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// const cityGameChan = "894280981098430514" // // REAL
const cityGameChan = "893415494512680990" // TEST

var City string

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	nick := getMessageAuthorNick(m)

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch {
	case strings.HasPrefix(m.Content, "!roll"), strings.HasPrefix(m.Content, "!ролл"):
		roll := game_roll(m.Content, nick)
		s.ChannelMessageSend(m.ChannelID, roll)

	case m.Content == "!bottle", m.Content == "!бутылочка":
		randMember, err := game_bottle(s, m)
		if err == nil {
			s.ChannelMessageSend(m.ChannelID, nick+" :kiss: "+getMemberNick(randMember))
		}

	case m.Content == "!score", m.Content == "!очки":
		score, err := b.repository.GetScore(m.Author.ID)
		if err == nil {
			s.ChannelMessageSend(m.ChannelID, nick+" your score is "+score)
		}

	case m.ChannelID == cityGameChan && (strings.HasPrefix(m.Content, "!city ") || strings.HasPrefix(m.Content, "!город ") || strings.HasPrefix(m.Content, "г ")):
		cityGame := game_cities(b, m)
		s.ChannelMessageSend(m.ChannelID, cityGame)
	}
}
