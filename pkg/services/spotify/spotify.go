package spotify_services

import (
	"io"
	"os"
	"log"
	"fmt"
	"net/http"
	"encoding/json"
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
func getTrackInfo(client *http.Client, trackID string) (map[string]interface{}, error) {
	// Send the GET request
	resp, err := http.Get(fmt.Sprintf("%s/tracks/%s", SPOTIFY_URL, trackID))
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON into a map
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}

func CreateSpotifyPlaylist(playlistName string, trackIDs []string) error {
	// Create the request body
	body := map[string]interface{}{
		"name":        playlistName,
		"public":      false,
		"description": "Playlist created by OMI",
		"tracks":      trackIDs,
	}
	fmt.Println(body)
	return nil
}

// func AddTracksToSpotifyPlaylist(client *http.Client, playlistID, trackID string) error {
// 	url := "https://api.spotify.com/v1/users/" + "{user_id}" + "/playlists"
// 	client.Post(url, "application/json")
// 	// Create the request body
// 	body := map[string]interface{}{
// 		"uris": []string{fmt.Sprintf("spotify:track:%s", trackID)},
// 	}
// 	fmt.Println("Uris:", body)
// 	return nil
// }

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
