package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
		denyChan, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "duelRoll.denyChan",
		})
		sendMessageReply(s, m, denyChan)
		return
	}

	str := strings.Fields(m.Content)

	// command validation
	if len(str) != 3 {
		wrongFormat, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "duelRoll.wrongFormat",
		})
		sendMessageReply(s, m, wrongFormat)
		return
	}
	if len(m.Mentions) != 1 {
		opponentNotFound, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "duelRoll.opponentNotFound",
		})
		sendMessageReply(s, m, opponentNotFound)
		return
	}
	bet, err := strconv.Atoi(str[2])
	if err != nil {
		shouldBeNumber, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "bet.shouldBeNumber",
		})
		sendMessageReply(s, m, shouldBeNumber)
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
		denyChan, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "duelRoll.denyChan",
		})
		sendRespond(s, i, denyChan)
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
		denyBot, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "duelRoll.denyBot",
		})
		return denyBot, false
	}

	if duel.bet <= 0 {
		greaterThanZero, _ := h.localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "bet.greaterThanZero",
		})
		return greaterThanZero, false
	}

	authorScore, _ := h.repository.GetScore(duel.challenger.ID)
	opponentScore, _ := h.repository.GetScore(duel.opponent.ID)
	if duel.bet > authorScore {
		betTooHigh := h.localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID:   "bet.betTooHigh",
			PluralCount: authorScore,
		})
		return betTooHigh, false
	}
	if duel.bet > opponentScore {
		betTooHighForOpponent := h.localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID:   "duelRoll.betTooHighForOpponent",
			PluralCount: opponentScore,
		})
		return betTooHighForOpponent, false
	}

	duelInvite, _ := h.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "duelRoll.duelInvite",
		TemplateData: map[string]string{
			"Opponent": duel.opponent.Mention(),
		},
	})
	return duelInvite, true
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

	challengerRoll := rand.Intn(100) + 1 //nolint:gosec
	opponentRoll := rand.Intn(100) + 1   //nolint:gosec

	var sb strings.Builder
	challengerRollMsg, _ := h.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "duelRoll.rollMsg",
		TemplateData: map[string]string{
			"Name": duel.challenger.Username,
			"Roll": strconv.Itoa(challengerRoll),
		},
	})
	opponentRollMsg, _ := h.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "duelRoll.rollMsg",
		TemplateData: map[string]string{
			"Name": duel.opponent.Username,
			"Roll": strconv.Itoa(opponentRoll),
		},
	})
	challengerWinMsg, _ := h.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "duelRoll.winMsg",
		TemplateData: map[string]string{
			"Name":  duel.challenger.Username,
			"Score": strconv.Itoa(duel.bet),
		},
		PluralCount: duel.bet,
	})
	opponentWinMsg, _ := h.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "duelRoll.winMsg",
		TemplateData: map[string]string{
			"Name":  duel.opponent.Username,
			"Score": strconv.Itoa(duel.bet),
		},
		PluralCount: duel.bet,
	})
	equal, _ := h.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "duelRoll.equal",
	})

	sb.WriteString(challengerRollMsg)
	sb.WriteString(opponentRollMsg)
	switch {
	case challengerRoll < opponentRoll:
		h.updateScore(duel.opponent.ID, duel.challenger.ID, duel.bet)
		sb.WriteString(opponentWinMsg)
	case challengerRoll > opponentRoll:
		h.updateScore(duel.challenger.ID, duel.opponent.ID, duel.bet)
		sb.WriteString(challengerWinMsg)
	default:
		sb.WriteString(equal)
	}
	return sb.String(), nil
}

func (h *Handler) updateScore(winnerID, loserID string, bet int) {
	err := h.repository.AddScore(winnerID, bet)
	if err != nil {
		log.Printf("error during increasing score: %v\n", err)
		return
	}
	err = h.repository.AddScore(loserID, -bet)
	if err != nil {
		_ = h.repository.AddScore(winnerID, -bet)
		log.Printf("error during decreasing score: %v\n", err)
	}
}
