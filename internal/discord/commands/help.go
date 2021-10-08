package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) help(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpCommand := "Команды:\n" +
		"- `!roll (!ролл)` Спорим 100 не выкинешь?\n" +
		"- `!bet (!бет, !ставка)` Лудомания\n" +
		"- `!bottle (!бутылочка)` Целуйтесь или бан!\n" +
		"- `!city (!город, !г)` Игра в города\n" +
		"- `!top (!топ, !leaderboard)` Список лидеров\n" +
		"- `!help (!помощь)` это сообщение"

	msg := discordgo.MessageEmbed{
		Description: helpCommand,
		Color:       0x006969, // 96 96
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Println(err)
	}
}
