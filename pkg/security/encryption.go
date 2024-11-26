package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func EncryptToken(token *oauth2.Token) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	// Convert the token to JSON
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("failed to serialize token: %v", err)
	}

	// Create a cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// Use GCM (Galois/Counter Mode) for encryption
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// Generate a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	encrypted := aesGCM.Seal(nonce, nonce, tokenJSON, nil)

	// Return the encrypted data as a base64 string
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func DecryptToken(encryptedToken string) (*oauth2.Token, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	// Decode the base64-encoded token
	cipherText, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %v", err)
	}

	// Create a cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	// Use GCM for decryption
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	// Extract the nonce size and decrypt
	nonceSize := aesGCM.NonceSize()
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %v", err)
	}

	// Convert the JSON back to a Token
	var token oauth2.Token
	if err := json.Unmarshal(plainText, &token); err != nil {
		return nil, fmt.Errorf("failed to deserialize token: %v", err)
	}

	return &token, nil
}
