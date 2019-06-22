package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/rylio/ytdl"
)

type Song2Upload struct {
	SongName   string
	Artist     string
	URL        string
	StartTime  string
	EndTime    string
	Lyric      string
	AlbumCover *string `json:"albumCover,omitempty"`
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

	resp, _ := http.Head(song.URL)
	vid, err := ytdl.GetVideoInfo(resp.Request.URL)
	if err != nil {
		return err
	}

	mp4File, _ := os.Create(shortID + ".mp4")
	defer mp4File.Close()
	fmt.Println(vid.Formats)
	err = vid.Download(vid.Formats[0], mp4File)
	if err != nil {
		fmt.Println("1")
		return err
	}

	// parse when to begin video
	startTime, err := strconv.Atoi(song.StartTime)
	if err != nil {
		fmt.Println("2")
		return err
	}

	// parse when to end video
	endTime, err := strconv.Atoi(song.EndTime)
	if err != nil {
		fmt.Println("3")
		return err
	}

	amount := strconv.Itoa(endTime - startTime)

	// slice video
	sliceFileName := shortID + ".mp4"
	_, err = exec.Command("ffmpeg", "-ss", song.StartTime, "-t", amount, "-i", mp4File.Name(), sliceFileName).Output()
	if err != nil {
		fmt.Println("4")
		return err
	}

	// extract audio from video
	rapbitFileName := shortID + ".mp3"
	_, err = exec.Command("ffmpeg", "-i", sliceFileName, "-f", "mp3", "-ab", "100000", "-vn", rapbitFileName).Output()
	if err != nil {
		fmt.Println("5")
		return err
	}

	// retriev album cover image
	if song.AlbumCover == nil {
		*song.AlbumCover = vid.GetThumbnailURL(ytdl.ThumbnailQualityDefault).String()
	}
	response, err := http.Get(*song.AlbumCover)
	if err != nil {
		fmt.Println("6")
		return err
	}
	defer response.Body.Close()

	// open a file for writing
	contentType := response.Header.Get("Content-Type")
	extType, err := mime.ExtensionsByType(contentType)
	if err != nil && len(extType) > 0 {
		fmt.Println("7")
		return err
	}

	albumCoverImageFile, err := os.Create(shortID + extType[0])
	if err != nil {
		fmt.Println("WORLD")
		return err
	}
	defer albumCoverImageFile.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(albumCoverImageFile, response.Body)
	if err != nil {
		fmt.Println("HELlo")
		return err
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
