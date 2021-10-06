package discord

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func game_bet(b *Bot, m *discordgo.MessageCreate) (string, error) {
	var scoreSign = ":tamale:"
	var newScore int
	var sb strings.Builder
	sb.WriteString(getMessageAuthorNick(m))

	str := strings.Split(m.Content, " ")
	if len(str) < 2 {
		sb.WriteString(" укажи ставку")
		return sb.String(), nil
	}
	bet, err := strconv.Atoi(str[1])
	if err != nil {
		return "", err
	}
	score, err := b.repository.GetScore(m.Author.ID)
	if err != nil {
		return "", err
	}
	scoreForOutput := strconv.Itoa(score)
	rand.Seed(time.Now().Unix())

	if bet > score {
		sb.WriteString(" ставка не может превышать количество ")
		sb.WriteString(scoreSign)
		sb.WriteString(" Всего у тебя ")
		sb.WriteString(scoreForOutput)
		sb.WriteString(scoreSign)
		return sb.String(), nil
	}

	roll := rand.Intn(100)

	sb.WriteString(" сделал ставку ")
	sb.WriteString(str[1])
	sb.WriteString(scoreSign)
	if roll < 51 {
		sb.WriteString(" и проиграл! :stuck_out_tongue_closed_eyes: ")
		b.repository.AddScore(m.Author.ID, -bet)
		newScore = score - bet

	} else {
		sb.WriteString(" и выйграл! :partying_face: ")
		b.repository.AddScore(m.Author.ID, bet)
		newScore = score + bet
	}
	sb.WriteString(" Теперь у тебя ")
	sb.WriteString(strconv.Itoa(newScore))
	sb.WriteString(scoreSign)

	return sb.String(), nil
}
