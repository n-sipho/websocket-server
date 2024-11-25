package recognition

import (
	"github.com/acrcloud/acrcloud_sdk_golang/acrcloud"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func RecognizeAudio(filePath string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
