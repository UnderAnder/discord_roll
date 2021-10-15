package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
)

// scoreMessage Print user score to the guild channel in response to the text command
func (h *Handler) scoreMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	score, err := h.score(m.Author.ID)
	if err != nil {
		return
	}
	if _, err := s.ChannelMessageSend(m.ChannelID, m.Author.Username+" у тебя "+score+" очков"); err != nil {
		log.Printf("Failed to response the command %v, %v\n", m.Content, err)
	}
}

// scoreSlash Print user score to the guild channel in response to the slash command
func (h *Handler) scoreSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		return
	}
	score, err := h.score(i.Member.User.ID)
	if err != nil {
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: i.Member.User.Username + " у тебя " + score + " очков",
		},
	})
	if err != nil {
		log.Printf("Failed to response the command %v, %v\n", i.ApplicationCommandData().Name, err)
	}
}

// score Return user score from DB as string
func (h *Handler) score(userID string) (string, error) {
	score, err := h.repository.GetScore(userID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strconv.Itoa(score), nil
}
