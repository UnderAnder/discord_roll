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

type emptyStruct struct{}

var (
	void         emptyStruct
	prevCity     string
	prevAuthorID string
	prevTime     time.Time
	prevCities   = make(map[string]struct{})
)

// cityMessage Output the result of the cities game in response to the text command
func (h *Handler) cityMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != cityGameChan && m.ChannelID != cityGameChanTest {
		return
	}

	str := strings.SplitN(m.Content, " ", 2)
	city := strings.ToLower(str[1])

	result := h.city(m.Author.ID, city)

	if _, err := s.ChannelMessageSendReply(m.ChannelID, result, m.Message.Reference()); err != nil {
		log.Printf("Failed to response the command %v, %v\n", m.Content, err)
	}
}

// citySlash Output the result of the cities game on the guild channel in response to the slash command
func (h *Handler) citySlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := interactionUserID(i)

	city := i.ApplicationCommandData().Options[0].StringValue()
	text := h.city(userID, city)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
		},
	})
	if err != nil {
		log.Printf("Failed to response the command %v, %v\n", i.ApplicationCommandData().Name, err)
	}
}

// city Return result of a cities game as string
func (h *Handler) city(discordID, city string) string {
	var sb strings.Builder
	cityTitle := strings.Title(city)
	exists, _ := h.repository.CityExist(city)
	username := discordgo.User{ID: discordID}.Username

	sb.WriteString(username)

	_, alreadyGuessed := prevCities[city]
	if alreadyGuessed {
		sb.WriteString(" город ")
		sb.WriteString(cityTitle)
		sb.WriteString(" уже был назван")
		return sb.String()
	}

	if !exists {
		sb.WriteString(" город ")
		sb.WriteString(cityTitle)
		sb.WriteString(" не существует")
		return sb.String()
	}

	// start game
	if prevCity == "" {
		prevCity = city
		prevAuthorID = discordID
		prevTime = time.Now()
		prevCities[city] = void
		sb.WriteString(" Игра началась, следующий город на ")
		sb.WriteString(strings.ToUpper(getLastChar(prevCity)))
		return sb.String()
	}

	lastChar := getLastChar(prevCity)

	if strings.HasPrefix(city, lastChar) {
		score := scoreAccrual(discordID)
		err := h.repository.AddScore(discordID, score)
		if err != nil {
			log.Printf("Failed to change score for userID: %v, %v\n", discordID, err)
		}

		prevCity = city
		prevAuthorID = discordID
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
