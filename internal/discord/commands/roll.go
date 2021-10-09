package commands

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) roll(s *discordgo.Session, m *discordgo.MessageCreate) {
	str := strings.Fields(m.Content)
	maxScore := 100
	quantity := 1

	if len(str) == 4 && str[1] == "duel" {
		h.duelRoll(s, m, str)
		return
	}

	switch len(str) {
	case 1:
		break
	case 3:
		var err error
		quantity, err = strconv.Atoi(str[2])
		if err != nil {
			quantity = 1
		}
		fallthrough
	case 2:
		var err error
		maxScore, err = strconv.Atoi(str[1])
		if err != nil {
			maxScore = 100
		}
	}

	// build output string
	var sb strings.Builder
	sb.WriteString(getMessageAuthorNick(m))
	sb.WriteString(" (1-")
	sb.WriteString(strconv.Itoa(maxScore))
	sb.WriteString(") ")
	for i := 0; i < quantity; i++ {
		sb.WriteString(" :game_die:")
		sb.WriteString(strconv.Itoa(rand.Intn(maxScore) + 1)) //nolint:gosec
	}

	if _, err := s.ChannelMessageSend(m.ChannelID, sb.String()); err != nil {
		log.Println(err)
	}
}

func (h *Handler) duelRoll(s *discordgo.Session, m *discordgo.MessageCreate, str []string) {
	var sb strings.Builder
	if len(m.Mentions) != 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Оппонент не найден, упомяни участника через @")
		if err != nil {
			return
		}
		return
	}

	opponent := m.Mentions[0]

	bet, err := strconv.Atoi(str[3])
	if err != nil {
		log.Println(err)
		return
	}
	if bet < 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ставка не может быть отрицательной")
		return
	}

	authorScore, _ := h.repository.GetScore(m.Author.ID)
	opponentScore, _ := h.repository.GetScore(opponent.ID)
	if bet > authorScore {
		sb.WriteString("Ставка слишком высока, у тебя всего ")
		sb.WriteString(strconv.Itoa(authorScore))
		if _, err := s.ChannelMessageSend(m.ChannelID, sb.String()); err != nil {
			log.Println(err)
		}
		return
	}
	if bet > opponentScore {
		sb.WriteString("У твоего оппонента недостаточно очков, ставка не должна превышать ")
		sb.WriteString(strconv.Itoa(opponentScore))
		if _, err := s.ChannelMessageSend(m.ChannelID, sb.String()); err != nil {
			log.Println(err)
		}
		return
	}

	sb.WriteString(opponent.Mention())
	sb.WriteString(" тебя вызвали на дуэль, нажми на :game_die: чтобы принять, или :no_entry_sign: чтобы отказаться")
	message, err := s.ChannelMessageSend(m.ChannelID, sb.String())
	if err != nil {
		log.Println(err)
	}
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
	if err != nil {
		return
	}
	err = s.MessageReactionAdd(message.ChannelID, message.ID, "🎲")
	if err != nil {
		return
	}
	err = s.MessageReactionAdd(message.ChannelID, message.ID, "🚫")
	if err != nil {
		return
	}

	accept := <-h.eventChan
	if accept != "roll" {
		_, _ = s.ChannelMessageEdit(message.ChannelID, message.ID, opponent.Username+" отказался")
		return
	}

	authorRoll := rand.Intn(100) + 1   //nolint:gosec
	opponentRoll := rand.Intn(100) + 1 //nolint:gosec
	sb.Reset()
	sb.WriteString(getMessageAuthorNick(m) + " выбрасывает " + strconv.Itoa(authorRoll) + "\n")
	sb.WriteString(opponent.Username + " выбрасывает " + strconv.Itoa(opponentRoll) + "\n")
	if authorRoll == opponentRoll {
		sb.WriteString("Ваши силы равны, оба остались при своём")
	} else if authorRoll > opponentRoll {
		err := h.repository.AddScore(m.Author.ID, bet)
		if err != nil {
			return
		}
		err = h.repository.AddScore(opponent.ID, -bet)
		if err != nil {
			return
		}
		sb.WriteString(getMessageAuthorNick(m) + " победил и получает " + strconv.Itoa(bet) + " очков соперника!")
	} else {
		err := h.repository.AddScore(m.Author.ID, -bet)
		if err != nil {
			return
		}
		err = h.repository.AddScore(opponent.ID, bet)
		if err != nil {
			return
		}
		sb.WriteString(opponent.Username + " победил и получает" + strconv.Itoa(bet) + " очков соперника!")
	}

	_, err = s.ChannelMessageSend(m.ChannelID, sb.String())
	if err != nil {
		log.Println(err)
	}
}
