package spotify_services

import (
	"encoding/json"
	"fmt"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"os"
	"websocket-server/pkg/database"
)

var (
	SPOTIFY_URL = os.Getenv("SPOTIFY_BASE_URL")
)

type user struct {
	Name string `json:"display_name"`
	ID   string `json:"id"`
}

// getTrackInfo retrieves information about a Spotify track by its ID.
// It sends a GET request to the Spotify API and parses the JSON response
// into a map[string]interface{}.
//
// trackID is the ID of the Spotify track to retrieve information for.
// Returns the track information as a map, or an error if the request failed.
// func getTrackInfo(client *http.Client, trackID string) (map[string]interface{}, error) {
// 	// Send the GET request
// 	resp, err := http.Get(fmt.Sprintf("%s/tracks/%s", SPOTIFY_URL, trackID))
// 	if err != nil {
// 		return nil, fmt.Errorf("error making GET request: %w", err)
// 	}
// 	defer resp.Body.Close()
// 	// Read the response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response body: %w", err)
// 	}

// 	// Parse the JSON into a map
// 	var result map[string]interface{}
// 	err = json.Unmarshal(body, &result)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
// 	}

// 	return result, nil
// }

func CreateSpotifyPlaylist(client *spotify.Client, uid string) (spotify.ID, error) {
	userId, err := database.GetUser(uid)
	if err != nil {
		log.Println("Error getting user from database:", err)
		return "", err
	}

	// Fetch user's playlists
	playlists, err := client.GetPlaylistsForUser(userId)
	if err != nil {
		log.Println("Error fetching user playlists:", err)
		return "", err
	}

	// Check if the playlist already exists
	var existingPlaylistID spotify.ID
	for _, playlist := range playlists.Playlists {
		if playlist.Name == "OMI Songs" {
			existingPlaylistID = playlist.ID
			break
		}
	}

	// If the playlist exists, return its ID
	if existingPlaylistID != "" {
		return existingPlaylistID, nil
	}

	// Create a new playlist if it doesn't exist
	playlist, err := client.CreatePlaylistForUser(userId, "OMI Songs", "Playlist created by OMI", false)
	if err != nil {
		log.Println("Error creating playlist:", err)
		return "", err
	}

	return playlist.ID, nil
}

func AddTrackToSpotifyPlaylist(client *spotify.Client, trackId spotify.ID, playlistId spotify.ID) error {
	playlistTracks, getPlaylistError := client.GetPlaylistTracks(playlistId)
	if getPlaylistError != nil {
		log.Printf("Error getting playlist tracks: %v", getPlaylistError)
		return getPlaylistError
	}
	// Check if the track is already in the playlist
	for _, track := range playlistTracks.Tracks {
		if track.Track.ID == trackId {
			log.Println("Song already exists in the playlist.")
			return nil // Song already exists, no need to add
		}
	}

	// Add the track to the playlist
	_, err := client.AddTracksToPlaylist(playlistId, trackId)
	if err != nil {
		log.Printf("Error adding tracks: %v", err)
	}
	log.Printf("Added %v track to playlist", trackId)
	return nil
}

func GetSpotifyUserInfo(client *http.Client) (*user, error) {
	resp, err := client.Get("https://api.spotify.com/v1/me")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from Spotify: %s", resp.Status)
	}

	var userInfo user
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	log.Printf("User ID: %s", userInfo.ID)
	return &userInfo, nil
}
