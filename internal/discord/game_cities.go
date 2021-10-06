package discord

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

type void struct{}

var (
	member       void
	prevCity     string
	prevAuthorID string
	prevTime     time.Time
	prevCities   = make(map[string]struct{})
)

func game_cities(b *Bot, m *discordgo.MessageCreate) string {
	str := strings.SplitN(m.Content, " ", 2)
	city := strings.ToLower(str[1])
	cityForOutput := strings.ToTitle(city)
	exists, _ := b.repository.CityExist(city)
	var sb strings.Builder

	sb.WriteString(getMessageAuthorNick(m))

	_, alredyGuessed := prevCities[city]
	if alredyGuessed {
		sb.WriteString(" город ")
		sb.WriteString(cityForOutput)
		sb.WriteString(" уже был назван")
		return sb.String()
	}

	if !exists {
		sb.WriteString(" город ")
		sb.WriteString(cityForOutput)
		sb.WriteString(" не существует")
		return sb.String()
	}

	if prevCity == "" {
		prevCity = city
		prevAuthorID = m.Author.ID
		prevTime = time.Now()
		prevCities[city] = member
		sb.WriteString(" Игра началась, следующий город на ")
		sb.WriteString(strings.ToUpper(getLastChar(prevCity)))
		return sb.String()
	}

	lastChar := getLastChar(prevCity)

	if strings.HasPrefix(city, lastChar) {
		score := scoreAccrual(m.Author.ID)
		b.repository.AddScore(m.Author.ID, score)

		prevCity = city
		prevAuthorID = m.Author.ID
		prevTime = time.Now()
		lastChar = getLastChar(prevCity)

		sb.WriteString("  :tada: +")
		sb.WriteString(strconv.Itoa(score))
		sb.WriteString(" Слудующий город на ")
	} else {
		sb.WriteString(" город должен начинаться на ")
	}
	sb.WriteString(strings.ToUpper(lastChar))

	return sb.String()
}

func scoreAccrual(id string) int {
	var score int

	switch id {
	case prevAuthorID:
		score += 1
	default:
		score += 6
	}

	switch {
	case time.Since(prevTime) < 3000000000:
		score += 6
	case time.Since(prevTime) < 5000000000:
		score += 4
	case time.Since(prevTime) < 10000000000:
		score += 3
	case time.Since(prevTime) < 15000000000:
		score += 2
	default:
		score += 1
	}

	return score
}

func getLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return strings.ToLower(s[len(s)-size:])
}
