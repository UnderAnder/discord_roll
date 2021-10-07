package repository

import "github.com/UnderAnder/discord_roll/internal/repository/sqlite"

type Repository interface {
	AddScore(discordID string, score int) error
	GetScore(discordID string) (int, error)
	CityExist(city string) (bool, error)
	GetTopUsersByScore(limit int) ([]sqlite.User, error)
}
