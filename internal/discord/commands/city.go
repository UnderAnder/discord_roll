package commands

import (
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

const cityGameChan = "894280981098430514"
const cityGameChanTest = "893415494512680990"

type void struct{}

var (
	member       void
	prevCity     string
	prevAuthorID string
	prevTime     time.Time
	prevCities   = make(map[string]struct{})
)

func (h *Handler) city(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != cityGameChan || m.ChannelID != cityGameChanTest {
		return
	}

	str := strings.SplitN(m.Content, " ", 2)
	city := strings.ToLower(str[1])
	cityForOutput := strings.ToTitle(city)
	exists, _ := h.repository.CityExist(city)
	var sb strings.Builder

	sb.WriteString(getMessageAuthorNick(m))

	_, alreadyGuessed := prevCities[city]
	if alreadyGuessed {
		sb.WriteString(" город ")
		sb.WriteString(cityForOutput)
		sb.WriteString(" уже был назван")
		_, err := s.ChannelMessageSend(m.ChannelID, sb.String())
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	if !exists {
		sb.WriteString(" город ")
		sb.WriteString(cityForOutput)
		sb.WriteString(" не существует")
		_, err := s.ChannelMessageSend(m.ChannelID, sb.String())
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	if prevCity == "" {
		prevCity = city
		prevAuthorID = m.Author.ID
		prevTime = time.Now()
		prevCities[city] = member
		sb.WriteString(" Игра началась, следующий город на ")
		sb.WriteString(strings.ToUpper(getLastChar(prevCity)))
		_, err := s.ChannelMessageSend(m.ChannelID, sb.String())
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	lastChar := getLastChar(prevCity)

	if strings.HasPrefix(city, lastChar) {
		score := scoreAccrual(m.Author.ID)
		err := h.repository.AddScore(m.Author.ID, score)
		if err != nil {
			log.Println(err)
			return
		}

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

	if _, err := s.ChannelMessageSend(m.ChannelID, sb.String()); err != nil {
		log.Println(err)
	}
}

func scoreAccrual(id string) int {
	var score int

	switch id {
	case prevAuthorID:
		score++
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
		score++
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
