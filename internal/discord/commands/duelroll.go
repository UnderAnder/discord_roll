package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type duel struct {
	sess       *discordgo.Session
	channelID  string
	messageID  string
	challenger *discordgo.User
	opponent   *discordgo.User
	bet        int
}

// duelMessage Print invite to a duel and the result of the duel to the guild channel in response to the text command
func (h *Handler) duelMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// check message came from a channel
	if m.Member == nil {
		sendMessageReply(s, m, "Вызвать на дуэль можно только на канале")
		return
	}

	str := strings.Fields(m.Content)

	// check is command correct
	if len(str) != 3 {
		sendMessageReply(s, m, "Неверный формат, должно быть `!duel @user ставка`")
		return
	}
	if len(m.Mentions) != 1 {
		sendMessageReply(s, m, "Оппонент не найден, упомяни одного участника через @")
		return
	}

	bet, err := strconv.Atoi(str[2])
	if err != nil {
		sendMessageReply(s, m, "Ставка должна быть числом")
		return
	}

	opponent := m.Mentions[0]
	duel := duel{sess: s, channelID: m.ChannelID, challenger: m.Author, opponent: opponent, bet: bet}
	// print invite to duel
	startTxt, ok := h.duelStart(duel)
	startMsg, done := sendMessageReply(s, m, startTxt)
	if !ok || !done {
		return
	}

	duel.messageID = startMsg.ID
	duelResult, err := h.duel(duel)
	if err != nil {
		log.Printf("Duel failed %v\n", err)
		return
	}
	if duelResult == "" {
		return
	}
	// update message to show result
	_, err = s.ChannelMessageEdit(m.ChannelID, startMsg.ID, duelResult)
	if err != nil {
		log.Println(err)
	}
}

// duelSlash Print invite to a duel and the result of the duel to the guild channel in response to the slash command
func (h *Handler) duelSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		sendRespond(s, i, "Вызвать на дуэль можно только на канале")
		return
	}

	opponent := i.ApplicationCommandData().Options[0].UserValue(s)
	bet := int(i.ApplicationCommandData().Options[1].IntValue())

	duel := duel{sess: s, channelID: i.ChannelID, challenger: i.Member.User, opponent: opponent, bet: bet}
	// print invite to duel
	startTxt, ok := h.duelStart(duel)
	done := sendRespond(s, i, startTxt)
	if !ok || !done {
		return
	}

	// get interaction message
	resp, err := s.InteractionResponse(s.State.User.ID, i.Interaction)
	if err != nil {
		log.Println(err)
		return
	}

	duel.messageID = resp.ID
	duelResult, err := h.duel(duel)
	if err != nil {
		log.Println(err)
		return
	}
	// there is nothing to print
	if duelResult == "" {
		return
	}
	// update message to show result
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
func (h *Handler) duelStart(duel duel) (string, bool) {
	// check opponent is not bot
	if duel.opponent.ID == duel.sess.State.User.ID {
		return "Нельзя вызвать бота на дуэль", false
	}

	if duel.bet < 0 {
		return "Ставка не может быть отрицательной", false
	}

	authorScore, _ := h.repository.GetScore(duel.challenger.ID)
	opponentScore, _ := h.repository.GetScore(duel.opponent.ID)
	if duel.bet > authorScore {
		return fmt.Sprintf("Ставка слишком высока, у тебя всего %d", authorScore), false
	}
	if duel.bet > opponentScore {
		msg := "У твоего оппонента недостаточно очков, ставка не должна превышать %d"
		return fmt.Sprintf(msg, opponentScore), false
	}

	msg := "%s тебя вызвали на дуэль, нажми на :game_die: чтобы принять, или :no_entry_sign: чтобы отказаться"
	return fmt.Sprintf(msg, duel.opponent.Mention()), true
}

// duel process the duel
func (h *Handler) duel(duel duel) (string, error) {
	err := duel.sess.MessageReactionAdd(duel.channelID, duel.messageID, "🎲")
	if err != nil {
		return "", err
	}
	err = duel.sess.MessageReactionAdd(duel.channelID, duel.messageID, "🚫")
	if err != nil {
		log.Println(err)
		// not critical no need for return here
	}

	accept := <-h.eventChan
	if accept != "roll" {
		_, _ = duel.sess.ChannelMessageEdit(duel.channelID, duel.messageID, duel.opponent.Username+" отказался")
		return "", nil
	}

	authorRoll := rand.Intn(100) + 1   //nolint:gosec
	opponentRoll := rand.Intn(100) + 1 //nolint:gosec

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s выбрасывает %d\n", duel.challenger.Username, authorRoll))
	sb.WriteString(fmt.Sprintf("%s выбрасывает %d\n", duel.opponent.Username, opponentRoll))
	switch {
	case authorRoll < opponentRoll:
		_ = h.repository.AddScore(duel.challenger.ID, -duel.bet)
		_ = h.repository.AddScore(duel.opponent.ID, duel.bet)
		sb.WriteString(fmt.Sprintf("%s победил и получает %d очков соперника!", duel.opponent.Username, duel.bet))
	case authorRoll > opponentRoll:
		_ = h.repository.AddScore(duel.challenger.ID, duel.bet)
		_ = h.repository.AddScore(duel.opponent.ID, -duel.bet)

		sb.WriteString(fmt.Sprintf("%s победил и получает %d очков соперника!", duel.challenger.Username, duel.bet))
	default:
		sb.WriteString("Ваши силы равны, оба остались при своём")
	}
	return sb.String(), nil
}
