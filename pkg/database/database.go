package database

import (
	"log"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// DB is a global variable for the SQLite database connection
var DB *sql.DB

func InitDB() {
	log.Println("Initializing database...")
	var err error
	DB, err = sql.Open("sqlite3", "./tracks.db") // Open a connection to the SQLite database file named app.db
	if err != nil {
		log.Fatal(err) // Log an error and stop the program if the database can't be opened
	}

	// SQL statement to create the todos table if it doesn't exist
	sqlStmt := `
			CREATE TABLE IF NOT EXISTS tracks (
    			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    			title VARCHAR(255) NOT NULL,
    			artist VARCHAR(255) NOT NULL,
    			spotify_id VARCHAR(255) NOT NULL UNIQUE,
    			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			);

			CREATE TABLE IF NOT EXISTS users (
    			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				uid VARCHAR(255) NOT NULL,
    			user_id VARCHAR(255) NOT NULL,
    			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			);
			`

	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt) // Log an error if table creation fails
	}
}

func AddUser(uid, spotifyUserId string) {
	log.Printf("Adding user with uid: %s and spotifyUserId: %s", uid, spotifyUserId)
	_, err := DB.Exec("INSERT INTO users(uid, user_id) VALUES(?, ?)", uid, spotifyUserId)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return
	}
	log.Printf("User added to database with uid: %s and spotifyUserId: %s", uid, spotifyUserId)
}

func GetUser(uid string) (string, error) {
	var spotifyUserId string
	err := DB.QueryRow("SELECT user_id FROM users WHERE uid = ?", uid).Scan(&spotifyUserId)
	if err != nil {
		return "", err
	}
	return spotifyUserId, nil
}

// saveTrack handles the creation of a new trackId from spotify in the database
func SaveTrack(title, artist, spotifyId string) {
	log.Println("Checking for existing track with spotify_id:", spotifyId)

	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tracks WHERE spotify_id = ?)", spotifyId).Scan(&exists)
	if err != nil {
		log.Printf("Error checking for existing track: %v", err)
		return
	}

	if exists {
		log.Printf("Track with spotify_id %s already exists, skipping save", spotifyId)
		return
	}

	log.Println("Saving track to database with title: ", title, " and artist: ", artist)
	createdAt := time.Now()
	_, err = DB.Exec("INSERT INTO tracks(title, artist, spotify_id, created_at) VALUES(?, ?, ?, ?)", title, artist, spotifyId, createdAt) // Insert the new todo into the database
	if err != nil {
		log.Printf("Error inserting track into database: %v", err)
		return
	}
	log.Println("Track saved to database with title: ", title, " and artist: ", artist)
}

type TrackId struct {
	ID string
}

func GetAllTracks() ([]TrackId, error) {
	rows, err := DB.Query("SELECT spotify_id FROM tracks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []TrackId
	for rows.Next() {
		var track TrackId
		if err := rows.Scan(&track.ID); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tracks, nil

}
