package sqlite

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

type Sqlite struct {
	db       *sql.DB
	location string
}

type User struct {
	DiscordID string
	Score     int
}

// NewSqlite builds and returns a pointer to a Sqlite repository implementation
func NewSqlite(location string) (*Sqlite, error) {
	return &Sqlite{location: location}, nil
}

// Close closes the connection to sqlite.
func (s *Sqlite) Close() error {
	log.Println("Closing sqlite db...")
	if s.db == nil {
		log.Println("DB already closed, nothing to do.")
		return nil
	}
	if err := s.db.Close(); err != nil {
		return err
	}
	log.Println("Closed sqlite db.")
	return nil
}

// Open opens a connection to sqlite.
func (s *Sqlite) Open() error {
	log.Printf("Opening sqlite db at '%s'...", s.location)
	db, err := sql.Open("sqlite3", s.location)
	if err != nil {
		return err
	}
	s.db = db
	log.Println("Opened sqlite db.")
	return nil
}

func (s *Sqlite) AddScore(discordID string, score int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into users(discord_id, score) values(?, ?) ON CONFLICT(discord_id) DO UPDATE SET score=score+?;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(discordID, score, score)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite) GetScore(discordID string) (int, error) {
	stmt, err := s.db.Prepare("select score from users where discord_id = ?")
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

func (s *Sqlite) CityExist(c string) (bool, error) {
	stmt, err := s.db.Prepare("select title_ru from cities where title_ru = ?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	var city string
	err = stmt.QueryRow(strings.Title(c)).Scan(&city)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func (s *Sqlite) GetTopUsersByScore(limit int) ([]User, error) {
	stmt := `select discord_id, score from users order by score desc limit ?`

	rows, err := s.db.Query(stmt, limit)
	if err != nil {
		log.Println(err)
		return []User{}, err
	}
	defer rows.Close()

	var result []User

	for rows.Next() {
		item := User{}
		err := rows.Scan(&item.DiscordID, &item.Score)
		if err != nil {
			log.Println(err)
			return []User{}, err
		}
		result = append(result, item)
	}
	return result, nil
}
