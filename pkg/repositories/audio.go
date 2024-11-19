package repositories

import (
	"log"
	"time"
	"websocket-server/pkg/database"
)

// saveTrack handles the creation of a new trackId from spotify in the database
func SaveTrack(title, artist, spotifyId string) {
	DB := database.DB
	log.Println("Saving track to database with title: ", title, " and artist: ", artist)
	// Insert the new track into the database
	createdAt := time.Now()
	_, err := DB.Exec("INSERT INTO tracks(title, artist, spotify_id, created_at) VALUES(?, ?, ?, ?)", title, artist, spotifyId, createdAt) // Insert the new todo into the database
	if err != nil {
		log.Printf("Error inserting track into database: %v", err)
		return
	}
	log.Println("Track saved to database with title: ", title, " and artist: ", artist)
}
