package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"websocket-server/pkg/database"
	"websocket-server/pkg/handlers"
)

func main() {
	port := "8080"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// oauthService := NewOAuthService()

	// API endpoints
	http.HandleFunc("/login/spotify", handlers.HandleSpotifyLogin)
	http.HandleFunc("/callback/spotify", handlers.HandleCallback)
	http.HandleFunc("/audio", handlers.HandlePostAudio)

	// Database setup
	database.InitDB()
	defer database.DB.Close()
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
