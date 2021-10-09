package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
)

func (h *Handler) score(s *discordgo.Session, m *discordgo.MessageCreate) {
	score, err := h.repository.GetScore(m.Author.ID)
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := s.ChannelMessageSend(m.ChannelID, getMessageAuthorNick(m)+" у тебя "+strconv.Itoa(score)+" очков"); err != nil {
		log.Println(err)
	}
}
