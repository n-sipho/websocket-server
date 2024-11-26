package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"websocket-server/pkg/audio"
	"websocket-server/pkg/database"
	"websocket-server/pkg/models/song_recognition_response"
	"websocket-server/pkg/services/recognition"
	spotify_services "websocket-server/pkg/services/spotify"

	"github.com/zmb3/spotify"

	// "websocket-server/pkg/services/s"

	"websocket-server/pkg/utils"
)

const redirectURI = "http://localhost:8080/callback/spotify"

var (
	auth = spotify.NewAuthenticator(
		redirectURI,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPrivate,
		spotify.ScopePlaylistModifyPublic,
	)
	ch    = make(chan *spotify.Client)
	state = utils.GenerateRandomState()
)

func HandlePostAudio(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	uid := query.Get("uid")
	omiUserId := strings.Split(uid, "?")[0]

	log.Printf("Received request from uid: %s", omiUserId)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	currentTime := time.Now()
	filename := fmt.Sprintf("%02d_%02d_%04d_%02d_%02d_%02d",
		currentTime.Day(),
		currentTime.Month(),
		currentTime.Year(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second())

	tempFilePath := filepath.Join(os.TempDir(), filename)

	header := audio.CreateWAVHeader(len(body))

	// Write to temporary file
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Printf("Failed to create temp file: %v", err)
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Write WAV header and audio data
	tempFile.Write(header)
	tempFile.Write(body)

	// cloudUloadResponse, err := uploadAudioToCLD(filename, tempFilePath)
	// if err != nil {
	// 	log.Printf("Failed to upload audio: %v", err)
	// 	http.Error(w, "Failed to upload audio", http.StatusInternalServerError)
	// 	return
	// }

	// Recognize audio
	log.Printf("Song Recognition started for audio file: %s", tempFilePath)
	songRecognitionRes := recognition.RecognizeAudio(tempFilePath)
	var response songrecognitionresponse.SongRecognitionResponse

	// Parse the JSON string
	err = json.Unmarshal([]byte(songRecognitionRes), &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// client := <-ch
	// client.AddTracksToPlaylist()

	if len(response.Metadata.Music) > 0 {
		// save the response to the database
		title := response.Metadata.Music[0].Title
		artist := response.Metadata.Music[0].Artists[0].Name
		spotifyID := response.Metadata.Music[0].ExternalMetadata.Spotify.Track.ID

		database.SaveTrack(title, artist, spotifyID)

		token, tokenError := database.GetSpotifyToken(omiUserId)
		if tokenError != nil {
			log.Printf("Failed to get token: %v", tokenError)
			http.Error(w, "Failed to get token", http.StatusInternalServerError)
			return
		}

		client := auth.NewClient(token)
		playListId, createPlaylistError := spotify_services.CreateSpotifyPlaylist(&client, omiUserId)
		if createPlaylistError != nil {
			log.Printf("Failed to create playlist: %v", createPlaylistError)
			http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
			return
		}

		// Add track to playlist
		addTrackError := spotify_services.AddTrackToSpotifyPlaylist(&client, spotifyID, playListId)
		if addTrackError != nil {
			log.Printf("Failed to add track to playlist: %v", addTrackError)
			http.Error(w, "Failed to add track to playlist", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

}

func HandleSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	url := auth.AuthURL(state)
	// Redirect user to Spotify auth page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	// uid := r.URL.Query().Get("uid") // For production
	uid := "dlCfOXqkmgMqSr29ycTgzohl4gp1" // For testing purposes
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := auth.NewClient(tok)
	user, err := client.CurrentUser()
	if err != nil {
		http.Error(w, "Couldn't get user info", http.StatusInternalServerError)
		return
	}

	dbUser, dbUserError := database.GetUser(uid)
	if dbUserError != nil && dbUser == "" {
		database.AddUser(uid, user.ID)
	}

	saveTokenError := database.SaveSpotifyToken(tok, uid)
	if saveTokenError != nil {
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		return
	}

	response := map[string]bool{
		"is_setup_completed": true,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
