package game

import (
	"math/rand"
	"strconv"
	"strings"
)

func Roll(message, nick string) string {
	str := strings.Split(message, " ")
	maxScore := 100
	quantity := 1
	var sb strings.Builder

	if len(str) > 1 {
		var err error
		maxScore, err = strconv.Atoi(str[1])
		if err != nil {
			maxScore = 100
		}
		if len(str) == 3 {
			quantity, err = strconv.Atoi(str[2])
			if err != nil {
				quantity = 1
			}
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
