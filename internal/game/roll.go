package game

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func Roll(message, nick string) string {
	rand.Seed(time.Now().Unix())

	str := strings.Split(message, " ")
	maxScore := 100
	quantity := 1
	var sb strings.Builder

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
	sb.WriteString(nick)
	sb.WriteString(" (1-")
	sb.WriteString(strconv.Itoa(maxScore))
	sb.WriteString(") ")
	for i := 0; i < quantity; i++ {
		sb.WriteString(" :game_die:")
		sb.WriteString(strconv.Itoa(rand.Intn(maxScore) + 1))
	}

	return sb.String()
}
