package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	SPOTIFY_URL = os.Getenv("SPOTIFY_BASE_URL")
)

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

func AddTrackToSpotifyPlaylist(playlistID, trackID string) error {
	// Create the request body
	body := map[string]interface{}{
		"uris": []string{fmt.Sprintf("spotify:track:%s", trackID)},
	}
	fmt.Println("Uris:", body)
	return nil
}
