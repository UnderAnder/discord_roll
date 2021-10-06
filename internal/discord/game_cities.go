package discord

import (
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

var (
	prevCity     string
	prevAuthorID string
)

func game_cities(b *Bot, m *discordgo.MessageCreate) string {
	str := strings.SplitN(m.Content, " ", 2)
	city := str[1]
	exist, _ := b.repository.CityExist(city)
	var sb strings.Builder

	sb.WriteString(getMessageAuthorNick(m))

	if !exist {
		sb.WriteString(city)
		sb.WriteString(" не существует")
		return sb.String()
	}

	if prevCity == "" {
		prevCity = city
		prevAuthorID = m.Author.ID
		sb.WriteString(" Игра началась, следующий город на ")
		sb.WriteString(strings.ToUpper(getLastChar(prevCity)))
		return sb.String()
	}

	lastChar := getLastChar(prevCity)

	if strings.HasPrefix(city, strings.ToLower(lastChar)) {
		prevCity = city
		lastChar = getLastChar(prevCity)

		sb.WriteString("Верно! Слудующий город на ")
	} else {
		sb.WriteString(" город должен начинаться на ")
	}
	sb.WriteString(strings.ToUpper(lastChar))

	return sb.String()
}

func getLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return strings.ToLower(s[len(s)-size:])
}
