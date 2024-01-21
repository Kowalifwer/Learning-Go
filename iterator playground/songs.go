package main

import "fmt"

// Song represents a music track
type Song struct {
	Title  string
	Artist string
}

func main() {
	// Create a playlist
	songs := []Song{
		{"Wonderwall", "Oasis"},
		{"Don't Look Back in Anger", "Oasis"},
		{"Champagne Supernova", "Oasis"},
		{"Slide Away", "Oasis"},
	}

	var playlist []Song = songs

	// Iterate over the songs in the playlist using range
	for i, song := range playlist {
		fmt.Printf("%d. Now playing: %s by %s\n", i+1, song.Title, song.Artist)
	}
}
