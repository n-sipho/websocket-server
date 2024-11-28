package recognition

import (
	// "github.com/acrcloud/acrcloud_sdk_golang/acrcloud"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	// "os"
	"strconv"
	"time"
)

type ACRCloudOptions struct {
	Host             string
	Endpoint         string
	SignatureVersion string
	DataType         string
	AccessKey        string
	AccessSecret     string
}

func buildStringToSign(method, uri, accessKey, dataType, signatureVersion string, timestamp int64) string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%d", method, uri, accessKey, dataType, signatureVersion, timestamp)
}

func signString(signString, accessSecret string) string {
	mac := hmac.New(sha1.New, []byte(accessSecret))
	mac.Write([]byte(signString))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func RecognizeAudio(filePath string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("ACRCLOUD_HOST")
	accessKey := os.Getenv("ACRCLOUD_ACCESS_KEY")
	accessSecret := os.Getenv("ACRCLOUD_SECRET_ACCESS_KEY")

	options := ACRCloudOptions{
		Host:             host,
		Endpoint:         "/v1/identify",
		SignatureVersion: "1",
		DataType:         "audio",
		AccessKey:        accessKey,
		AccessSecret:     accessSecret,
	}
	// Open the file (this could be a temporary file or a regular file)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get file info to get the size
	fi, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fi.Size()

	// Read the file into a buffer
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("access_key", options.AccessKey)
	writer.WriteField("data_type", options.DataType)
	writer.WriteField("signature_version", options.SignatureVersion)

	// Generate the timestamp and signature
	timestamp := time.Now().Unix()
	writer.WriteField("timestamp", strconv.FormatInt(timestamp, 10))

	// Sign the request
	stringToSign := buildStringToSign("POST", options.Endpoint, options.AccessKey, options.DataType, options.SignatureVersion, timestamp)
	signature := signString(stringToSign, options.AccessSecret)
	writer.WriteField("signature", signature)

	// Append the file to the form
	writer.WriteField("sample_bytes", strconv.FormatInt(fileSize, 10))

	// Create the form file field for the audio sample
	sampleWriter, err := writer.CreateFormFile("sample", filePath)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(sampleWriter, file)
	if err != nil {
		return "", err
	}

	// Close the writer so we can get the final form data
	err = writer.Close()
	if err != nil {
		return "", err
	}

	// Send the request
	url := fmt.Sprintf("http://%s%s", options.Host, options.Endpoint)
	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body completely
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
