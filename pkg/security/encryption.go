package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
	"crypto/rand"

)

func EncryptToken(plainText string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("ENCRYPTION_KEY")
	// Convert key and plaintext to byte slices
	keyBytes := []byte(key)
	plainTextBytes := []byte(plainText)

	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Use GCM for encryption
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a nonce (unique number used once per encryption)
	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	// Encrypt the plaintext
	cipherText := aesGCM.Seal(nonce, nonce, plainTextBytes, nil)

	// Return the encrypted text as a Base64-encoded string
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptToken(cipherText, key string) (string, error) {
	// Convert key and ciphertext to byte slices
	keyBytes := []byte(key)
	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Use GCM for decryption
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract the nonce from the ciphertext
	nonceSize := aesGCM.NonceSize()
	if len(cipherTextBytes) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, cipherTextBytes := cipherTextBytes[:nonceSize], cipherTextBytes[nonceSize:]

	// Decrypt the ciphertext
	plainTextBytes, err := aesGCM.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}

	// Return the decrypted text as a string
	return string(plainTextBytes), nil
}

