package main

import (
	"log"
	"net/http"
)

func main() {

	InitDB()

	http.HandleFunc("/upload", uploadRoute)
	http.HandleFunc("/song/", songRoute)
	http.HandleFunc("/songs", songsRoute)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
