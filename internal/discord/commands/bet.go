package commands

import (
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
	if len(str) != 2 {
		result = "Укажи ставку"
	} else {
		bet, err := strconv.Atoi(str[1])
		if err != nil {
			log.Println(err)
			result = "Ставка должна быть числом"
		} else {
			result = h.bet(m.Author.ID, bet)
		}
	}

	if _, err := s.ChannelMessageSendReply(m.ChannelID, result, m.Message.Reference()); err != nil {
		log.Printf("Failed to response the command %v, %v\n", m.Content, err)
	}
}

// betSlash Print result of a bet in response to the slash command
func (h *Handler) betSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := interactionUserID(i)

	bet := int(i.ApplicationCommandData().Options[0].IntValue())
	text := h.bet(userID, bet)

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

// bet Return result of a bet as string
func (h *Handler) bet(discordID string, bet int) string {
	var scoreSign = " очков"
	var newScore int
	var sb strings.Builder

	score, err := h.repository.GetScore(discordID)
	if err != nil {
		log.Println(err)
		sb.WriteString(" у тебя нет очков, чтобы совершить ставку")
		return sb.String()
	}

	scoreStr := strconv.Itoa(score)

	if bet > score {
		sb.WriteString(" Слишком высокая ставка, у тебя всего ")
		sb.WriteString(scoreStr)
		sb.WriteString(scoreSign)
		return sb.String()
	}

	roll := rand.Intn(100) //nolint:gosec

	if roll < 52 {
		sb.WriteString("Проиграл! :stuck_out_tongue_closed_eyes: ")
		err := h.repository.AddScore(discordID, -bet)
		if err != nil {
			log.Printf("Failed to change score for userID: %v, %v\n", discordID, err)
		}
		newScore = score - bet
	} else {
		sb.WriteString("Выйграл! :partying_face: ")
		err := h.repository.AddScore(discordID, bet)
		if err != nil {
			log.Printf("Failed to change score for userID: %v, %v\n", discordID, err)
		}
		newScore = score + bet
	}
	sb.WriteString(" Теперь у тебя ")
	sb.WriteString(strconv.Itoa(newScore))
	sb.WriteString(scoreSign)
	return sb.String()
}
