package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

// duelMessage Print invite to a duel and the result of the duel to the guild channel in response to the text command
func (h *Handler) duelMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	str := strings.Fields(m.Content)
	if len(str) != 3 {
		return
	}
	if len(m.Mentions) != 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Оппонент не найден, упомяни одного участника через @")
		if err != nil {
			log.Println(err)
		}
		return
	}

	bet, err := strconv.Atoi(str[2])
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, "Ставка должна быть числом")
		if err != nil {
			log.Println(err)
		}
		return
	}

	opponent := m.Mentions[0]
	startTxt := h.duelStart(m.Author, opponent, bet)

	startMsg, err := s.ChannelMessageSend(m.ChannelID, startTxt)
	if err != nil {
		log.Println(err)
		return
	}
	duelResult, err := h.duel(s, m.ChannelID, startMsg.ID, m.Author, opponent, bet)
	if err != nil {
		log.Println(err)
		return
	}
	if duelResult == "" {
		return
	}
	_, err = s.ChannelMessageEdit(m.ChannelID, startMsg.ID, duelResult)
	if err != nil {
		log.Println(err)
	}
}

// duelSlash Print invite to a duel and the result of the duel to the guild channel in response to the slash command
func (h *Handler) duelSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		return
	}

	opponent := i.ApplicationCommandData().Options[0].UserValue(s)
	bet := int(i.ApplicationCommandData().Options[1].IntValue())

	startTxt := h.duelStart(i.Member.User, opponent, bet)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: startTxt,
		},
	})
	if err != nil {
		log.Printf("Failed to response the command %v, %v\n", i.ApplicationCommandData().Name, err)
	}
	resp, _ := s.InteractionResponse(s.State.User.ID, i.Interaction)

	duelResult, err := h.duel(s, i.ChannelID, resp.ID, i.Member.User, opponent, bet)
	if err != nil {
		log.Println(err)
		return
	}
	if duelResult == "" {
		return
	}
	_, err = s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
		Content: duelResult,
	})
	if err != nil {
		_, _ = s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
		return
	}
}

// duelStart invite to a duel
func (h *Handler) duelStart(challenger, opponent *discordgo.User, bet int) string {
	var sb strings.Builder

	if bet < 0 {
		return "Ставка не может быть отрицательной"
	}

	authorScore, _ := h.repository.GetScore(challenger.ID)
	opponentScore, _ := h.repository.GetScore(opponent.ID)
	if bet > authorScore {
		sb.WriteString("Ставка слишком высока, у тебя всего ")
		sb.WriteString(strconv.Itoa(authorScore))
		return sb.String()
	}
	if bet > opponentScore {
		sb.WriteString("У твоего оппонента недостаточно очков, ставка не должна превышать ")
		sb.WriteString(strconv.Itoa(opponentScore))
		return sb.String()
	}

	sb.WriteString(opponent.Mention())
	sb.WriteString(" тебя вызвали на дуэль, нажми на :game_die: чтобы принять, или :no_entry_sign: чтобы отказаться")
	return sb.String()
}

// duel process the duel
func (h *Handler) duel(s *discordgo.Session, channelID, messageID string, challenger, opponent *discordgo.User, bet int) (string, error) {
	betStr := strconv.Itoa(bet)
	err := s.MessageReactionAdd(channelID, messageID, "🎲")
	if err != nil {
		return "", err
	}
	err = s.MessageReactionAdd(channelID, messageID, "🚫")
	if err != nil {
		log.Println(err)
		// not critical no need for return here
	}

	accept := <-h.eventChan
	if accept != "roll" {
		_, _ = s.ChannelMessageEdit(channelID, messageID, opponent.Username+" отказался")
		return "", nil
	}

	authorRoll := rand.Intn(100) + 1   //nolint:gosec
	opponentRoll := rand.Intn(100) + 1 //nolint:gosec

	var sb strings.Builder
	sb.WriteString(challenger.Username + " выбрасывает " + strconv.Itoa(authorRoll) + "\n")
	sb.WriteString(opponent.Username + " выбрасывает " + strconv.Itoa(opponentRoll) + "\n")
	switch {
	case authorRoll < opponentRoll:
		err := h.repository.AddScore(challenger.ID, -bet)
		if err != nil {
			return "", err
		}
		err = h.repository.AddScore(opponent.ID, bet)
		if err != nil {
			return "", err
		}
		sb.WriteString(opponent.Username + " победил и получает" + betStr + " очков соперника!")
	case authorRoll > opponentRoll:
		err := h.repository.AddScore(challenger.ID, bet)
		if err != nil {
			return "", err
		}
		err = h.repository.AddScore(opponent.ID, -bet)
		if err != nil {
			return "", err
		}
		sb.WriteString(challenger.Username + " победил и получает " + betStr + " очков соперника!")
	default:
		sb.WriteString("Ваши силы равны, оба остались при своём")
	}
	return sb.String(), nil
}
