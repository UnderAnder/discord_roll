package commands

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// betMessage Print result of a bet in response to the text command
func (h *Handler) betMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	var result string
	str := strings.Fields(m.Content)

	// check bet
	switch len(str) {
	case 2:
		bet, err := strconv.Atoi(str[1])
		if err != nil {
			log.Println(err)
			result = "Ставка должна быть числом"
			break
		}
		result = h.bet(m.Author.ID, bet)
	default:
		result = "Укажи ставку"
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
	var scoreSign = "очков"
	var newScore int
	var sb strings.Builder

	if bet < 1 {
		return "Ставка должна быть больше 0"
	}

	score, err := h.repository.GetScore(discordID)
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("У тебя нет %s, чтобы совершить ставку", scoreSign)
	}

	if bet > score {
		return fmt.Sprintf("Слишком высокая ставка, у тебя всего %d %s", score, scoreSign)
	}

	roll := rand.Intn(100) //nolint:gosec

	switch {
	case roll <= 51:
		sb.WriteString("Проиграл! :stuck_out_tongue_closed_eyes: ")
		bet = -bet
	case roll > 51:
		sb.WriteString("Выйграл! :tada: ")
	}

	err = h.repository.AddScore(discordID, bet)
	if err != nil {
		log.Printf("Failed to change score for userID: %v, %v\n", discordID, err)
	}
	newScore = score + bet

	sb.WriteString(fmt.Sprintf(" Теперь у тебя %d %s", newScore, scoreSign))
	return sb.String()
}
