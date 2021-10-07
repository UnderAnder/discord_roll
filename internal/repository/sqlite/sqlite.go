package sqlite

import (
	"database/sql"
	"log"
	"strings"
)

type Repository struct {
	db *sql.DB
}

type User struct {
	Discord_id string
	Score      int
}

// constructor
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddScore(discordID string, score int) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println(err)
	}
	stmt, err := tx.Prepare("insert into users(discord_id, score) values(?, ?) ON CONFLICT(discord_id) DO UPDATE SET score=score+?;")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(discordID, score, score)
	if err != nil {
		log.Println(err)
		return err
	}
	tx.Commit()

	return nil
}

func (r *Repository) GetScore(discordID string) (int, error) {
	stmt, err := r.db.Prepare("select score from users where discord_id = ?")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()
	var score int
	err = stmt.QueryRow(discordID).Scan(&score)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return score, nil
}

func (r *Repository) CityExist(c string) (bool, error) {
	stmt, err := r.db.Prepare("select title_ru from cities where title_ru = ?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	var city string
	err = stmt.QueryRow(strings.Title(c)).Scan(&city)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			log.Fatal(err)
			return false, err
		}
	}
	return true, nil
}

func (r *Repository) GetTopUsersByScore(limit int) ([]User, error) {
	stmt := `select discord_id, score from users order by score desc limit ?`

	rows, err := r.db.Query(stmt, limit)
	if err != nil {
		log.Println(err)
		return []User{}, err
	}
	defer rows.Close()

	var result []User

	for rows.Next() {
		item := User{}
		err := rows.Scan(&item.Discord_id, &item.Score)
		if err != nil {
			log.Println(err)
			return []User{}, err
		}
		result = append(result, item)
	}
	return result, nil
}
