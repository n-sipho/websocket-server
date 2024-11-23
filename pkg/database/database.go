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
}

// saveTrack handles the creation of a new trackId from spotify in the database
func SaveTrack(title, artist, spotifyId string) {
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
