package commands

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) bet(s *discordgo.Session, m *discordgo.MessageCreate) {
	var scoreSign = ":tamale:"
	var newScore int
	var sb strings.Builder
	sb.WriteString(getMessageAuthorNick(m))

	str := strings.Split(m.Content, " ")
	if len(str) < 2 {
		sb.WriteString(" укажи ставку")
		_, err := s.ChannelMessageSend(m.ChannelID, sb.String())
		if err != nil {
			log.Println(err)
		}
		return
	}
	bet, err := strconv.Atoi(str[1])
	if err != nil {
		log.Println(err)
		return
	}
	score, err := h.repository.GetScore(m.Author.ID)
	if err != nil {
		log.Println(err)
		return
	}
	scoreForOutput := strconv.Itoa(score)

	if bet > score {
		sb.WriteString(" ставка не может превышать количество ")
		sb.WriteString(scoreSign)
		sb.WriteString(" Всего у тебя ")
		sb.WriteString(scoreForOutput)
		sb.WriteString(scoreSign)
		_, err := s.ChannelMessageSend(m.ChannelID, sb.String())
		if err != nil {
			log.Println(err)
		}
		return
	}

	roll := rand.Intn(100) //nolint:gosec

	sb.WriteString(" сделал ставку ")
	sb.WriteString(str[1])
	sb.WriteString(scoreSign)
	if roll < 52 {
		sb.WriteString(" и проиграл! :stuck_out_tongue_closed_eyes: ")
		err := h.repository.AddScore(m.Author.ID, -bet)
		if err != nil {
			return
		}
		newScore = score - bet
	} else {
		sb.WriteString(" и выйграл! :partying_face: ")
		err := h.repository.AddScore(m.Author.ID, bet)
		if err != nil {
			return
		}
		newScore = score + bet
	}
	sb.WriteString(" Теперь у тебя ")
	sb.WriteString(strconv.Itoa(newScore))
	sb.WriteString(scoreSign)

	if _, err := s.ChannelMessageSend(m.ChannelID, sb.String()); err != nil {
		log.Println(err)
	}
}
