package main

import (
	"io"
	"os"
	"log"
	"fmt"
	"time"
	"context"
	"net/http"
	"math/rand"
	"path/filepath"
	"encoding/json"
	"encoding/binary"
	"golang.org/x/oauth2"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/spotify"
	"websocket-server/pkg/database"
	"websocket-server/pkg/services/spotify"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/acrcloud/acrcloud_sdk_golang/acrcloud"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"websocket-server/pkg/models/song_recognition_response"
)

const (
	numChannels   = 1 // Mono audio
	sampleRate    = 16000
	bitsPerSample = 16 // 16 bits per sample
)

type OAuthService struct {
	config *oauth2.Config
	state  string
}

func generateRandomState() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Intn(100000))
}

func NewOAuthService() *OAuthService {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	spotifyClientId := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifySecret := os.Getenv("SPOTIFY_SECRET")
	return &OAuthService{
		config: &oauth2.Config{
			ClientID:     spotifyClientId,
			ClientSecret: spotifySecret,
			Endpoint:     spotify.Endpoint,
			RedirectURL:  "http://localhost:8080/callback/spotify",
			Scopes:       []string{"playlist-modify-private", "playlist-read-private"},
		},
		state: generateRandomState(), // Use a dynamic state in production
	}
}

func uploadAudioToCLD(fileName string, filePath string) (*uploader.UploadResult, error) {
	// Add your Cloudinary credentials, set configuration parameter
	// Secure=true to return "https" URLs, and create a context
	//===================
	cld, err := cloudinary.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}
	cld.Config.URL.Secure = true
	ctx := context.Background()

	// Upload the image.
	// Set the asset's public ID and allow overwriting the asset with new versions
	resp, err := cld.Upload.Upload(ctx, filePath, uploader.UploadParams{
		PublicID:       fileName,
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(true),
		AssetFolder:    "omi-audio-files",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Log the delivery URL
	fmt.Println("****2. Upload an image****\nDelivery URL:", resp.SecureURL)
	return resp, nil
}

// CreateWAVHeader generates a WAV header for the given data length
func createWAVHeader(dataLength int) []byte {
	byteRate := sampleRate * numChannels * bitsPerSample / 8
	blockAlign := numChannels * bitsPerSample / 8
	header := make([]byte, 44)

	copy(header[0:4], []byte("RIFF"))
	binary.LittleEndian.PutUint32(header[4:8], uint32(36+dataLength))
	copy(header[8:12], []byte("WAVE"))

	copy(header[12:16], []byte("fmt "))
	binary.LittleEndian.PutUint32(header[16:20], 16)
	binary.LittleEndian.PutUint16(header[20:22], 1)
	binary.LittleEndian.PutUint16(header[22:24], uint16(numChannels))
	binary.LittleEndian.PutUint32(header[24:28], uint32(sampleRate))
	binary.LittleEndian.PutUint32(header[28:32], uint32(byteRate))
	binary.LittleEndian.PutUint16(header[32:34], uint16(blockAlign))
	binary.LittleEndian.PutUint16(header[34:36], bitsPerSample)

	copy(header[36:40], []byte("data"))
	binary.LittleEndian.PutUint32(header[40:44], uint32(dataLength))

	return header
}

func recognizeAudio(filePath string) string {
	host := os.Getenv("ACRCLOUD_HOST")
	accessKey := os.Getenv("ACRCLOUD_ACCESS_KEY")
	accessSecret := os.Getenv("ACRCLOUD_SECRET_ACCESS_KEY")

	configs := map[string]string{
		"access_key":     accessKey,
		"access_secret":  accessSecret,
		"host":           host,
		"recognize_type": acrcloud.ACR_OPT_REC_AUDIO,
	}

	var recHandler = acrcloud.NewRecognizer(configs)

	userParams := map[string]string{}

	result := recHandler.RecognizeByFile(filePath, 0, 10, userParams)
	return result
}

func handlePostAudio(w http.ResponseWriter, r *http.Request) {
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

	header := createWAVHeader(len(body))

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
	songRecognitionRes := recognizeAudio(tempFilePath)
	var response songrecognitionresponse.SongRecognitionResponse

	// Parse the JSON string
	err = json.Unmarshal([]byte(songRecognitionRes), &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if len(response.Metadata.Music) > 0 {
		// save the response to the database
		title := response.Metadata.Music[0].Title
		artist := response.Metadata.Music[0].Artists[0].Name
		spotifyID := response.Metadata.Music[0].ExternalMetadata.Spotify.Track.ID

		database.SaveTrack(title, artist, spotifyID)
		w.WriteHeader(http.StatusOK)
	}
}

func (o *OAuthService) HandleSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	url := o.config.AuthCodeURL(o.state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (o *OAuthService) HandleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	uid := r.URL.Query().Get("uid")
	if state != o.state {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	token, err := o.config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		log.Println("Token exchange error:", err)
		return
	}

	client := o.config.Client(context.Background(), token)
	userInfo, err := spotify_services.GetSpotifyUserInfo(client)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		log.Println("User info error:", err)
		return
	}
	database.AddUser(uid, userInfo.ID)
	log.Printf("User logged in: %s", userInfo.Name)
	fmt.Fprintf(w, "Logged in successfully! User: %s\n", userInfo.Name)

	response := map[string]bool{
		"is_setup_completed": true,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	oauthService := NewOAuthService()
	port := "8080"

	// API endpoints
	http.HandleFunc("/login/spotify", oauthService.HandleSpotifyLogin)
	http.HandleFunc("/callback/spotify", oauthService.HandleCallback)
	http.HandleFunc("/audio", handlePostAudio)

	// Database setup
	database.InitDB()
	defer database.DB.Close()
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
