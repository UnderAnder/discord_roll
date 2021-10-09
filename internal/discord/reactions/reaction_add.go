package reactions

import (
	"github.com/UnderAnder/discord_roll/internal/repository"
	"github.com/bwmarrin/discordgo"
	"log"
)

var duelMsg *discordgo.Message
var opponent *discordgo.User

type handler func(*discordgo.Session, *discordgo.MessageReactionAdd)

type Handler struct {
	repository repository.Repository
	eventChan  chan string
}

func NewHandler(r repository.Repository, e chan string) *Handler {
	return &Handler{repository: r, eventChan: e}
}

func (h *Handler) HandleAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	log.Println("emojiID: " + r.Emoji.ID + "emojName: " + r.Emoji.Name + "emojUser: ")
	if r.MessageReaction.UserID == s.State.User.ID && r.MessageReaction.Emoji.Name == "👍" {
		duelMsg, _ = s.ChannelMessage(r.ChannelID, r.MessageID)
		opponent = duelMsg.Mentions[0]
		log.Println("DBUG: " + opponent.ID)
	}

	if opponent == nil {
		return
	}
	if r.Emoji.Name == "🚫" && r.MessageReaction.UserID == opponent.ID {
		log.Println("Отказ")
		h.eventChan <- "reject"
	}
	if r.Emoji.Name == "🎲" && r.MessageReaction.UserID == opponent.ID {
		log.Println("Запустить ролл")
		h.eventChan <- "roll"
	}
}
