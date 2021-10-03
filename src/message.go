package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// This function will be called every time a new message is created
// on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	rand.Seed(time.Now().Unix())

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "!roll") {
		str := strings.Split(m.Content, " ")
		maxScore := 100
		quantity := 1
		var sb strings.Builder

		if len(str) > 1 {
			var err error
			maxScore, err = strconv.Atoi(str[1])
			if err != nil {
				maxScore = 100
			}
			if len(str) == 3 {
				quantity, err = strconv.Atoi(str[2])
				if err != nil {
					quantity = 1
				}
			}
		}

		sb.WriteString(getMessageAuthorNick(m))
		for i := 0; i < quantity; i++ {
			sb.WriteString(" :game_die:")
			sb.WriteString(strconv.Itoa(rand.Intn(maxScore) + 1))
		}
		s.ChannelMessageSend(m.ChannelID, sb.String())
	}

	if m.Content == "!bottle" {
		members, err := s.GuildMembers(m.GuildID, "", 1000)
		if err != nil {
			fmt.Println("error can't get guild members", err)
			return
		}

		randMember := members[rand.Intn(len(members))]
		for randMember.User.ID == m.Author.ID {
			randMember = members[rand.Intn(len(members))]
		}

		s.ChannelMessageSend(m.ChannelID, getMessageAuthorNick(m)+" :kiss: "+getMemberNick(randMember))
	}
}
