package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"log"
)

// scoreMessage Print user score in response to the text command
func (h *Handler) scoreMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	result, err := h.score(m.Author.ID)
	if err != nil {
		return
	}
	sendMessageReply(s, m, result)
}

// scoreSlash Print user score in response to the slash command
func (h *Handler) scoreSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := interactionUserID(i)
	result, err := h.score(userID)
	if err != nil {
		return
	}

	sendRespond(s, i, result)
}

// score Return user score from DB as string
func (h *Handler) score(userID string) (string, error) {
	score, err := h.repository.GetScore(userID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	result, _ := h.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:   "score.result",
		PluralCount: score,
	})
	return result, nil
}
