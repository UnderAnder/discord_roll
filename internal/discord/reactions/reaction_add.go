package reactions

import (
	"github.com/UnderAnder/discord_roll/internal/repository"
	"github.com/bwmarrin/discordgo"
)

var duelMsg *discordgo.Message
var opponent *discordgo.User

type Handler struct {
	repository repository.Repository
	eventChan  chan string
}

func NewHandler(r repository.Repository, e chan string) *Handler {
	return &Handler{repository: r, eventChan: e}
}

func (h *Handler) HandleAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	h.rollDuel(s, r)
}

func (h *Handler) rollDuel(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.MessageReaction.UserID == s.State.User.ID && r.MessageReaction.Emoji.Name == "🎲" {
		duelMsg, _ = s.ChannelMessage(r.ChannelID, r.MessageID)
		// expire previous game
		if opponent != nil {
			h.eventChan <- "reject"
		}
		if len(duelMsg.Mentions) != 1 {
			return
		}
		opponent = duelMsg.Mentions[0]
	}
	if opponent == nil {
		return
	}
	if r.Emoji.Name == "🚫" && r.MessageReaction.UserID == opponent.ID {
		opponent = nil
		h.eventChan <- "reject"
	}
	if r.Emoji.Name == "🎲" && r.MessageReaction.UserID == opponent.ID {
		opponent = nil
		h.eventChan <- "roll"
	}
}
