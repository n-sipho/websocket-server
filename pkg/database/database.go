package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"websocket-server/pkg/security"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
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

			CREATE TABLE IF NOT EXISTS tokens (
    			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    			token TEXT NOT NULL UNIQUE
			);

			CREATE INDEX IF NOT EXISTS tokens_token_idx ON tokens(token);

			CREATE INDEX IF NOT EXISTS tokens_user_id_idx ON tokens(user_id);
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

// GetUser retrieves the Spotify user ID associated with the given UID from OMI.
// It queries the users table to find the user_id column value where the uid column matches the provided uid argument.
// If a matching user is found, the Spotify user ID is returned. Otherwise, a empty string and an error are returned.
func GetUser(uid string) (string, error) {
	log.Printf("Getting user with uid: %s", uid)
	var spotifyUserId string
	err := DB.QueryRow("SELECT user_id FROM users WHERE uid = ?", uid).Scan(&spotifyUserId)
	if err != nil {
		return "", err
	}
	return spotifyUserId, nil
}

// saveTrack handles the creation of a new trackId from spotify in the database
func SaveTrack(title, artist string, spotifyId spotify.ID) error {
	log.Println("Checking for existing track with spotify_id:", spotifyId)

	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tracks WHERE spotify_id = ?)", spotifyId).Scan(&exists)
	if err != nil {
		log.Printf("Error checking for existing track: %v", err)
		return fmt.Errorf("error checking existing track: %w", err)
	}

	if exists {
		log.Printf("Track with spotify_id %s already exists, skipping save", spotifyId)
		return fmt.Errorf("track with spotify_id %s already exists", spotifyId)
	}

	log.Println("Saving track to database with title: ", title, " and artist: ", artist)
	createdAt := time.Now()
	_, err = DB.Exec("INSERT INTO tracks(title, artist, spotify_id, created_at) VALUES(?, ?, ?, ?)", title, artist, spotifyId, createdAt)
	if err != nil {
		log.Printf("Error inserting track into database: %v", err)
		return fmt.Errorf("error inserting track: %w", err)
	}
	log.Println("Track saved to database with title: ", title, " and artist: ", artist)
	return nil
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

func SaveSpotifyToken(token *oauth2.Token, userID string) error {
	log.Printf("Saving Spotify token for user %s", userID)
	encryptedToken, err := security.EncryptToken(token)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`INSERT INTO tokens (user_id, token) VALUES ($1, $2)`, userID, encryptedToken)
	return err
}

func GetSpotifyToken(userID string) (*oauth2.Token, error) {
	var token string
	err := DB.QueryRow(`SELECT token FROM tokens WHERE user_id = $1`, userID).Scan(&token)
	if err != nil {
		return nil, err
	}
	decryptedToken, err := security.DecryptToken(token)
	if err != nil {
		return nil, err
	}
	return decryptedToken, nil
}
