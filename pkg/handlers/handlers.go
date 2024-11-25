package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/zmb3/spotify"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"websocket-server/pkg/audio"
	"websocket-server/pkg/database"
	"websocket-server/pkg/models/song_recognition_response"
	"websocket-server/pkg/services/recognition"

	"websocket-server/pkg/utils"
)

const redirectURI = "http://localhost:8080/callback/spotify"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate)
	ch    = make(chan *spotify.Client)
	state = utils.GenerateRandomState()
)

func HandlePostAudio(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	sampleRateParam := query.Get("sample_rate")
	uid := query.Get("uid")

	log.Printf("Received request from uid: %s", uid)
	log.Printf("Requested sample rate: %s", sampleRateParam)

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
		w.WriteHeader(http.StatusOK)
	}

}

func HandleSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	url := auth.AuthURL(state)
	// Redirect user to Spotify auth page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	log.Println("Token expiry:", tok.Expiry)
	// client := auth.NewClient(tok)
	// user, err := client.CurrentUser()
	if err != nil {
		http.Error(w, "Couldn't get user info", http.StatusInternalServerError)
		return
	}

	response := map[string]bool{
		"is_setup_completed": true,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
