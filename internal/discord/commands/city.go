package commands

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

type emptyStruct struct{}

type cityGame struct {
	prevCity     string
	prevAuthorID string
	prevTime     time.Time
	prevCities   map[string]struct{}
}

var games = make(map[string]*cityGame)
var void emptyStruct

func newCityGame() *cityGame {
	return &cityGame{prevCities: make(map[string]struct{})}
}

// cityMessage Output the result of the cities game in response to the text command
func (h *Handler) cityMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	str := strings.SplitN(m.Content, " ", 2)
	city := strings.ToLower(str[1])

	result := h.city(m.ChannelID, m.Author.ID, city)
	sendMessageReply(s, m, result)
}

// citySlash Output the result of the cities game on the guild channel in response to the slash command
func (h *Handler) citySlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := interactionUserID(i)

	city := strings.ToLower(i.ApplicationCommandData().Options[0].StringValue())
	result := h.city(i.ChannelID, userID, city)
	sendRespond(s, i, result)
}

// city Return result of a cities game as string
func (h *Handler) city(channelID, discordID, city string) string {
	// check channel is allowed for game
	if h.cfg.Bot.CityChanID != "" && channelID != h.cfg.Bot.CityChanID {
		denyChan, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "city.denyChan",
		})
		return denyChan
	}
	// check game already started on channel
	game, ok := games[channelID]
	if !ok {
		game = newCityGame()
		games[channelID] = game
	}

	cityTitle := strings.Title(city)

	_, alreadyGuessed := game.prevCities[city]
	if alreadyGuessed {
		guessedMsg, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "city.alreadyGuessed",
			TemplateData: map[string]string{
				"CityTitle": cityTitle,
			},
		})
		return guessedMsg
	}

	// start game
	if game.prevCity == "" {
		game.prevCity = city
		game.prevAuthorID = discordID
		game.prevTime = time.Now()
		game.prevCities[city] = void
		startMsg, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "city.startMsg",
			TemplateData: map[string]string{
				"StartLetter": strings.ToUpper(getLastChar(game.prevCity)),
			},
		})
		return startMsg
	}

	lastChar := getLastChar(game.prevCity)
	if !strings.HasPrefix(city, lastChar) {
		shouldStart, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "city.shouldStart",
			TemplateData: map[string]string{
				"StartLetter": strings.ToUpper(lastChar),
			},
		})
		return shouldStart
	}

	exists, _ := h.repository.CityExist(city)
	if !exists {
		notExists, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "city.notExists",
			TemplateData: map[string]string{
				"CityTitle": cityTitle,
			},
		})
		return notExists
	}

	score := scoreAccrual(game, discordID)
	err := h.repository.AddScore(discordID, score)
	if err != nil {
		log.Printf("Failed to change score for userID: %v, %v\n", discordID, err)
	}

	game.prevCity = city
	game.prevAuthorID = discordID
	game.prevTime = time.Now()
	lastChar = getLastChar(game.prevCity)

	winMsg, _ := h.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "city.winMsg",
		TemplateData: map[string]string{
			"StartLetter": strings.ToUpper(lastChar),
			"Score":       strconv.Itoa(score),
		},
		PluralCount: score,
	})
	return winMsg
}

func scoreAccrual(game *cityGame, id string) int {
	var score int

	switch id {
	case game.prevAuthorID:
		score++
	default:
		score += 6
	}

	secondsPast := time.Since(game.prevTime).Seconds()
	switch {
	case secondsPast < 5:
		score += 6
	case secondsPast < 10:
		score += 4
	case secondsPast < 15:
		score += 3
	case secondsPast < 20:
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
	lastChar := strings.ToLower(s[len(s)-size:])
	// there is no city starts with "ь"
	if lastChar == "ь" {
		lastChar = strings.ToLower(s[len(s)-size*2 : len(s)-size])
	}
	return lastChar
}
