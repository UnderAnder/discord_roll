package commands

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// betMessage Print result of a bet in response to the text command
func (h *Handler) betMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	var result string
	str := strings.Fields(m.Content)

	// validation
	switch len(str) {
	case 2:
		bet, err := strconv.Atoi(str[1])
		if err != nil {
			shouldBeNumber, _ := h.localizer.Localize(&i18n.LocalizeConfig{
				MessageID: "bet.shouldBeNumber",
			})
			result = shouldBeNumber
			break
		}

		result = h.bet(m.Author.ID, bet)
	default:
		specifyBet, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "bet.specifyBet",
		})
		result = specifyBet
	}

	sendMessageReply(s, m, result)
}

// betSlash Print result of a bet in response to the slash command
func (h *Handler) betSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := interactionUserID(i)
	bet := int(i.ApplicationCommandData().Options[0].IntValue())
	result := h.bet(userID, bet)
	sendRespond(s, i, result)
}

// bet Return result of a bet as string
func (h *Handler) bet(discordID string, bet int) string {
	var newScore int
	var sb strings.Builder

	if bet < 1 {
		greaterThanZero, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "bet.greaterThanZero",
		})
		return greaterThanZero
	}

	score, err := h.repository.GetScore(discordID)
	if err != nil {
		log.Println(err)
		internalErr, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "internalErr",
		})
		return internalErr
	}

	if bet > score {
		betTooHigh := h.localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID:   "bet.betTooHigh",
			PluralCount: score,
		})
		return betTooHigh
	}

	roll := rand.Intn(100) //nolint:gosec

	switch {
	case roll <= 51:
		lose, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "bet.lose",
		})
		sb.WriteString(lose)
		bet = -bet
	case roll > 51:
		win, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "bet.win",
		})
		sb.WriteString(win)
	}

	err = h.repository.AddScore(discordID, bet)
	if err != nil {
		log.Printf("Failed to change score for userID: %v, %v\n", discordID, err)
	}
	newScore = score + bet

	scoreAfter := h.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:   "bet.scoreAfter",
		PluralCount: newScore,
	})
	sb.WriteString(scoreAfter)
	return sb.String()
}
