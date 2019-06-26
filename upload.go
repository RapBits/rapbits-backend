package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Song2Upload struct {
	SongName   string
	Artist     string
	URL        string
	StartTime  string
	EndTime    string
	Lyric      string
	AlbumCover string
}

// RandomString will generate a random `n` character alphanumeric string
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[r.Intn(len(letter))]
	}
	return string(b)
}

func uploadYTMP3Content(song *Song2Upload) error {

	var err error

	shortID := RandomString(30)

	// parse when to begin video
	startTime, err := strconv.Atoi(song.StartTime)
	if err != nil {
		return err
	}

	// parse when to end video
	endTime, err := strconv.Atoi(song.EndTime)
	if err != nil {
		return err
	}

	amount := strconv.Itoa(endTime - startTime)

	// download video using script
	cmd := exec.Command("/Users/umayahabdennabi/Desktop/github/go/src/github.com/rapbits-backend/download.sh", shortID, song.URL, song.StartTime, amount, song.AlbumCover)
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, stderr.String())
	}

	return err

}

func uploadRoute(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Expected POST request for /upload", http.StatusBadRequest)
		return
	}

	song := new(Song2Upload)
	if err := json.NewDecoder(r.Body).Decode(song); err != nil {
		http.Error(w, "Can't read body", http.StatusBadRequest)
		return
	}

	m, err := json.Marshal(song)
	if err != nil {
		http.Error(w, "Error encoding song in upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = uploadYTMP3Content(song)
	if err != nil {
		http.Error(w, "Error processing upload request: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(m)
}
