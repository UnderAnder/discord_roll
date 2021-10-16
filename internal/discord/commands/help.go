package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) help(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpCommand := "Команды:\n" +
		"- `!roll (!ролл)` Спорим 100 не выкинешь?\n" +
		"- `!duel (!дуэль)` Ролл против соперника со ставкой\n" +
		"- `!bet (!бет, !ставка)` Лудомания\n" +
		"- `!city (!город, !г)` Игра в города\n" +
		"- `!top (!топ, !leaderboard)` Список лидеров\n" +
		"- `!score (!очки)` Сколько у тебя очков\n" +
		"- `!help (!помощь)` это сообщение"

	msg := discordgo.MessageEmbed{
		Description: helpCommand,
		Color:       0x006969, // 96 96
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Println(err)
	}
}
