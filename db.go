package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

type Song struct {
	Song       string
	Lyric      string
	Artist     string
	AlbumCover string
	Mp3        string
	Mp4        string
	Tags       string
	ID         string
	Index      int64
}

// InitDB will inititiate our sql connection session that will be used
// by the routes
func InitDB() {
	var err error
	db, err = sql.Open("mysql", "root:wvEQYyrmLepvezvFDJ2@/rapbits")

	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

// GetSongs will return `x` amount of songs starting from an offset `y` given
// an offset `x` and a limit `y`
func GetSongs(offset int, limit int) ([]*Song, error) {

	rows, err := db.Query("SELECT * FROM songs WHERE `index` >= ? limit ?", offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	songs := make([]*Song, 0)
	for rows.Next() {
		song := new(Song)
		err := rows.Scan(&song.Song, &song.Lyric, &song.Artist, &song.AlbumCover, &song.Mp3, &song.Mp4, &song.Tags, &song.ID, &song.Index)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

// GetSong will return the song with the given `id`
func GetSong(id string) (*Song, error) {

	row := db.QueryRow("SELECT * FROM songs WHERE id = ?", id)

	song := new(Song)
	err := row.Scan(&song.Song, &song.Lyric, &song.Artist, &song.AlbumCover, &song.Mp3, &song.Mp4, &song.Tags, &song.ID, &song.Index)

	if err != nil {
		return nil, err
	}

	return song, err
}
