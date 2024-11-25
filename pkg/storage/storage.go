package storage

import (
    "context"
    "fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadAudioToCLD(fileName string, filePath string) (*uploader.UploadResult, error) {
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
