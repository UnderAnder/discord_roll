package repository

type Repository interface {
	AddScore(discordID string, score int) error
	GetScore(discordID string) (int, error)
	CityExist(city string) (bool, error)
}
