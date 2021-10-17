package commands

import (
	"fmt"
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
	sendMessageReply(s, m, result)
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

	result := h.roll(maxScore, quantity)
	sendRespond(s, i, result)
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
	sb.WriteString(fmt.Sprintf("(1-%d)", maxRoll))
	for i := 0; i < quantity; i++ {
		sb.WriteString(" :game_die:")
		sb.WriteString(strconv.Itoa(rand.Intn(maxRoll) + 1)) //nolint:gosec
	}
	return sb.String()
}
