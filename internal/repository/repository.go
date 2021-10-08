package repository

import "github.com/UnderAnder/discord_roll/internal/repository/sqlite"

type Repository interface {
	Open() error
	Close() error

	AddScore(discordID string, score int) error
	GetScore(discordID string) (int, error)
	CityExist(city string) (bool, error)
	GetTopUsersByScore(limit int) ([]sqlite.User, error)
}

func New(location string) (Repository, error) {
	return sqlite.NewSqlite(location)
}
