package main

import (
	"log"
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
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "!roll") || strings.HasPrefix(m.Content, "!ролл") {
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

		// build output string
		sb.WriteString(getMessageAuthorNick(m))
		sb.WriteString(" (1-")
		sb.WriteString(strconv.Itoa(maxScore))
		sb.WriteString(") ")
		for i := 0; i < quantity; i++ {
			sb.WriteString(" :game_die:")
			sb.WriteString(strconv.Itoa(rand.Intn(maxScore) + 1))
		}

		s.ChannelMessageSend(m.ChannelID, sb.String())
	}

	if m.Content == "!bottle" || m.Content == "!бутылочка" {
		members, err := s.GuildMembers(m.GuildID, "", 1000)
		if err != nil {
			log.Println("error can't get guild members", err)
			return
		}

		randMember := members[rand.Intn(len(members))]
		for randMember.User.ID == m.Author.ID {
			randMember = members[rand.Intn(len(members))]
		}

		s.ChannelMessageSend(m.ChannelID, getMessageAuthorNick(m)+" :kiss: "+getMemberNick(randMember))
	}

	if m.Content == "!score" || m.Content == "!очки" {
		score := getScore(m.Author.ID)

		s.ChannelMessageSend(m.ChannelID, getMessageAuthorNick(m)+" your score is "+score)
	}

	if strings.HasPrefix(m.Content, "!city ") || strings.HasPrefix(m.Content, "!город ") || strings.HasPrefix(m.Content, "г ") {
		str := strings.SplitN(m.Content, " ", 2)
		city := str[1]
		exist := cityExist(city)
		var sb strings.Builder

		sb.WriteString(getMessageAuthorNick(m))
		if exist {
			if City == "" {
				City = city
				sb.WriteString(" Игра началась следующий город на ")
				sb.WriteString(strings.ToUpper(getLastChar(City)))
			} else {
				lastChar := getLastChar(City)
				if strings.HasPrefix(city, lastChar) || strings.HasPrefix(city, strings.ToUpper(lastChar)) {
					addScore(m.Author.ID, 10)
					City = city
					lastChar = getLastChar(City)
					sb.WriteString(" :tada: +10 очков")
					sb.WriteString(" всего уже ")
					sb.WriteString(getScore(m.Author.ID))
					sb.WriteString(" Слудующий город на ")
					sb.WriteString(strings.ToUpper(lastChar))
				} else {
					sb.WriteString(" город должен начинаться на ")
					sb.WriteString(strings.ToUpper(lastChar))
				}
			}
		} else {
			sb.WriteString(" город ")
			sb.WriteString(strings.Title(city))
			sb.WriteString(" не существует")
		}

		s.ChannelMessageSend(m.ChannelID, sb.String())
	}
}
