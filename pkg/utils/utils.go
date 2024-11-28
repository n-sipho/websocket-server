package utils

import (
	"fmt"
	"math/rand"
	"time"
	// "websocket-server/pkg/models/song_recognition_response"
)

func GenerateRandomState() string {
	// Create a new random source seeded with current time
	source := rand.NewSource(time.Now().UnixNano())

	// Create a new random number generator from the source
	r := rand.New(source)

	// Generate a random integer and format as string
	return fmt.Sprintf("%d", r.Intn(100000))
}

// func IsEmptySpotifyMeta(spotifyMeta songrecognitionresponse.SpotifyMeta) bool {
//     return len(spotifyMeta.Artists) == 0 &&
//            (spotifyMeta.Track == SpotifyTrack{}) &&
//            (spotifyMeta.Album == SpotifyAlbum{})
// }