package songrecognitionresponse

import "github.com/zmb3/spotify"

type SongRecognitionResponse struct {
	CostTime   float64  `json:"cost_time"`
	ResultType int      `json:"result_type"`
	Metadata   Metadata `json:"metadata"`
	Status     Status   `json:"status"`
}

type Metadata struct {
	Music        []Music `json:"music"`
	TimestampUTC string  `json:"timestamp_utc"`
}

type Music struct {
	Title             string       `json:"title"`
	Genres            []Genre      `json:"genres"`
	Album             Album        `json:"album"`
	DurationMs        int          `json:"duration_ms"`
	DbBeginTimeOffset int          `json:"db_begin_time_offset_ms"`
	DbEndTimeOffset   int          `json:"db_end_time_offset_ms"`
	SampleBeginTime   int          `json:"sample_begin_time_offset_ms"`
	SampleEndTime     int          `json:"sample_end_time_offset_ms"`
	PlayOffset        int          `json:"play_offset_ms"`
	Label             string       `json:"label"`
	Acrid             string       `json:"acrid"`
	ExternalIds       interface{}  `json:"external_ids"`
	ExternalMetadata  ExternalMeta `json:"external_metadata"`
	ResultFrom        int          `json:"result_from"`
	ReleaseDate       string       `json:"release_date"`
	Score             int          `json:"score"`
	Artists           []Artist     `json:"artists"`
}

type Genre struct {
	Name string `json:"name"`
}

type Album struct {
	Name string `json:"name"`
}

type Artist struct {
	Name string `json:"name"`
}

type ExternalMeta struct {
	Spotify SpotifyMeta `json:"spotify"`
}

type SpotifyMeta struct {
	Artists []SpotifyArtist `json:"artists"`
	Track   SpotifyTrack    `json:"track"`
	Album   SpotifyAlbum    `json:"album"`
}

type SpotifyArtist struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type SpotifyTrack struct {
	Name string `json:"name"`
	ID   spotify.ID `json:"id"`
}

type SpotifyAlbum struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Status struct {
	Version string `json:"version"`
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
}
