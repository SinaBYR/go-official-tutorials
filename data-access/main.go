package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	var db *sql.DB

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if err != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected")

	albums, err := albumsByArtist("John Coltrane", db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %+v\n", albums)

	album, err := albumByID(1, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %+v\n", album)

	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	}, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %+v\n", albID)
}

func albumsByArtist(name string, db *sql.DB) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func albumByID(id int64, db *sql.DB) (Album, error) {
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)

	var alb Album
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumByID %d: no such album", id)
		}
		return alb, fmt.Errorf("albumByID %d: %v", id, err)
	}

	return alb, nil
}

func addAlbum(alb Album, db *sql.DB) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO album (title, artist, price) VALUES (?, ?, ?)",
		alb.Title, alb.Artist, alb.Price,
	)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}
