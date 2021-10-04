package discord

import (
	"math/rand"
	"strings"
	"time"

	"github.com/UnderAnder/discord_roll/internal/game"
	"github.com/bwmarrin/discordgo"
)

var City string

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	nick := getMessageAuthorNick(m)
	rand.Seed(time.Now().Unix())

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "!roll") || strings.HasPrefix(m.Content, "!ролл") {
		roll := game.Roll(m.Content, nick)

		s.ChannelMessageSend(m.ChannelID, roll)
	}

	if m.Content == "!bottle" || m.Content == "!бутылочка" {
		randMember, err := game.Bottle(s, m)
		if err == nil {
			s.ChannelMessageSend(m.ChannelID, getMessageAuthorNick(m)+" :kiss: "+getMemberNick(randMember))
		}
	}

	if m.Content == "!score" || m.Content == "!очки" {
		score, err := b.repository.GetScore(m.Author.ID)
		if err == nil {
			s.ChannelMessageSend(m.ChannelID, getMessageAuthorNick(m)+" your score is "+score)
		}
	}

	if strings.HasPrefix(m.Content, "!city ") || strings.HasPrefix(m.Content, "!город ") || strings.HasPrefix(m.Content, "г ") {
		str := strings.SplitN(m.Content, " ", 2)
		city := str[1]
		exist, _ := b.repository.CityExist(city)
		if exist {
			cityGame, _ := game.City(getMessageAuthorNick(m), city)
			s.ChannelMessageSend(m.ChannelID, cityGame)
		} else {
			s.ChannelMessageSend(m.ChannelID, city+" такой город не существует")
		}
	}
}
