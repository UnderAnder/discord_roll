package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func addScore(discordID string, score int) {
	tx, err := DB.Begin()
	if err != nil {
		log.Println(err)
	}
	stmt, err := tx.Prepare("insert into users(discord_id, score) values(?, ?) ON CONFLICT(discord_id) DO UPDATE SET score=score+?;")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(discordID, score, score)
	if err != nil {
		log.Println(err)
	}
	tx.Commit()

}

func getScore(discordID string) string {
	stmt, err := DB.Prepare("select score from users where discord_id = ?")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	var score int
	err = stmt.QueryRow(discordID).Scan(&score)
	if err != nil {
		log.Println(err)
	}
	return strconv.Itoa(score)
}

func cityExist(c string) bool {
	stmt, err := DB.Prepare("select title_ru from cities where title_ru = ?")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	var city string
	err = stmt.QueryRow(strings.Title(c)).Scan(&city)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
			return false
		}
	}
	return true
}
