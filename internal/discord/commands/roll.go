package commands

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// rollMessage Output a random numbers in response to the text command
func (h *Handler) rollMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error
	// if 0 default values will use
	var maxRoll = 0
	var quantity = 0

	str := strings.Fields(m.Content)

	// convert parameters if exists
	switch len(str) {
	case 3:
		quantity, err = strconv.Atoi(str[2])
		if err != nil {
			return
		}
		fallthrough
	case 2:
		maxRoll, err = strconv.Atoi(str[1])
		if err != nil {
			return
		}
	default:
		break
	}

	result := h.roll(maxRoll, quantity)

	if _, err := s.ChannelMessageSendReply(m.ChannelID, result, m.Message.Reference()); err != nil {
		log.Printf("Failed to response the command %v, %v\n", m.Content, err)
	}
}

// rollSlash Output a random numbers in response to the slash command
func (h *Handler) rollSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// if 0 default values will use
	var maxScore = 0
	var quantity = 0

	// Here we need to convert raw interface{} value to wanted type.
	if len(i.ApplicationCommandData().Options) >= 1 {
		maxScore = int(i.ApplicationCommandData().Options[0].IntValue())
	}
	if len(i.ApplicationCommandData().Options) == 2 {
		quantity = int(i.ApplicationCommandData().Options[1].IntValue())
	}

	msg := h.roll(maxScore, quantity)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, we'll discuss them in "responses" part
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		log.Printf("Failed to response the command %v, %v\n", i.ApplicationCommandData().Name, err)
	}
}

// roll Return a random numbers as string
func (h *Handler) roll(maxRoll, quantity int) string {
	// default values
	if maxRoll == 0 {
		maxRoll = 100
	}
	if quantity == 0 {
		quantity = 1
	}

	// build output string
	var sb strings.Builder
	sb.WriteString(" (1-")
	sb.WriteString(strconv.Itoa(maxRoll))
	sb.WriteString(") ")
	for i := 0; i < quantity; i++ {
		sb.WriteString(" :game_die:")
		sb.WriteString(strconv.Itoa(rand.Intn(maxRoll) + 1)) //nolint:gosec
	}
	return sb.String()
}
