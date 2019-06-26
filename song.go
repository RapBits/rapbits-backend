package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func songRoute(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Expected GET request for /songs", http.StatusBadRequest)
	}

	p := strings.Split(r.URL.Path, "/")
	var songID string
	if len(p) == 3 {
		songID = p[2]
	} else {
		http.Error(w, "incorrect request "+r.URL.Path, http.StatusBadRequest)
		return
	}
	song, err := GetSong(songID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TOOD: log which song was retrieved

	m, err := json.Marshal(song)
	if err != nil {
		errStr := "Unable to find song: " + err.Error()
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(m)

}
