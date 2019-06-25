package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/pkg/errors"
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

	// Get mp4 video format
	var downloadURL *url.URL
	for _, v := range vid.Formats {
		if v.Extension == "mp4" {
			downloadURL, _ = vid.GetDownloadURL(v)
			break
		}
	}

	if err != nil {
		return err
	}

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

	// slice video
	sliceFileName := shortID + ".mp4"
	cmd := exec.Command("ffmpeg", "-i", downloadURL.String(), "-ss", song.StartTime, "-t", amount, sliceFileName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, stderr.String())
	}

	// extract audio from video
	rapbitFileName := shortID + ".mp3"
	_, err = exec.Command("ffmpeg", "-i", sliceFileName, "-f", "mp3", "-ab", "100000", "-vn", rapbitFileName).Output()
	if err != nil {
		log.Fatal("Failed to perform ffmpeg")
		return err
	}

	// retriev album cover image
	var response *http.Response
	if song.AlbumCover == nil {
		response, err = http.Get(vid.GetThumbnailURL(ytdl.ThumbnailQualityDefault).String())
	} else {
		response, err = http.Get(*song.AlbumCover)
	}

	if err != nil {
		log.Fatal("Failed to GET album cover")
		return err
	}
	defer response.Body.Close()

	// open a file for writing
	contentType := response.Header.Get("Content-Type")
	extType, err := mime.ExtensionsByType(contentType)
	if err != nil || len(extType) < 1 {
		log.Fatal("Invalid ext type", extType)
		return err
	}

	albumCoverImageFile, err := os.Create(shortID + extType[0])
	if err != nil {
		log.Fatal("Failed to create album cover file")
		return err
	}
	defer albumCoverImageFile.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(albumCoverImageFile, response.Body)
	if err != nil {
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
