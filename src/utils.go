package main

import (
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

// UTILS

func getMessageAuthorNick(m *discordgo.MessageCreate) string {
	if m.Member.Nick != "" {
		return m.Member.Nick
	} else {
		return m.Author.Username
	}
}

func getMemberNick(m *discordgo.Member) string {
	if m.Nick != "" {
		return m.Nick
	} else {
		return m.User.Username
	}
}

func getLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[len(s)-size:]
}
