package game

import (
	"strings"
	"unicode/utf8"
)

var prevCity string

func City(nick, city string) (string, error) {
	var sb strings.Builder

	sb.WriteString(nick)

	if prevCity == "" {
		prevCity = city
		sb.WriteString(" Игра началась следующий город на ")
		sb.WriteString(strings.ToUpper(getLastChar(prevCity)))
	} else {
		lastChar := getLastChar(prevCity)
		if strings.HasPrefix(city, lastChar) || strings.HasPrefix(city, strings.ToUpper(lastChar)) {
			prevCity = city
			lastChar = getLastChar(prevCity)

			sb.WriteString(" Слудующий город на ")
			sb.WriteString(strings.ToUpper(lastChar))
		} else {
			sb.WriteString(" город должен начинаться на ")
			sb.WriteString(strings.ToUpper(lastChar))
		}
	}

	return sb.String(), nil
}

func getLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[len(s)-size:]
}
