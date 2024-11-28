package spotify_services

import (
	"github.com/zmb3/spotify"
	"log"
	"websocket-server/pkg/database"
)


type user struct {
	Name string `json:"display_name"`
	ID   string `json:"id"`
}


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

