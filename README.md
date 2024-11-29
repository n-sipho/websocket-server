# OMI Audio Streaming Server

A Go-based audio streaming server that handles real-time audio processing, song recognition, and cloud storage integration.

## Features

- Real-time audio processing with WAV format support
- Song recognition using ACRCloud API
- Cloudinary integration for audio file storage
- RESTful API endpoint for audio processing
- Database integration for track storage
<!-- 
## Prerequisites

- Go 1.x
- ACRCloud account credentials
- Cloudinary account credentials
- PostgreSQL database

## Environment Variables

Create a `.env` file with the following variables:

ACR_ACCESS_KEY=your_acr_access_key
ACR_ACCESS_SECRET=your_acr_access_secret
ACR_HOST=your_acr_host
CLOUDINARY_CLOUD_NAME=your_cloudinary_cloud_name
CLOUDINARY_API_KEY=your_cloudinary_api_key
CLOUDINARY_API_SECRET=your_cloudinary_api_secret
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name

## Installation


git clone https://github.com/yourusername/omi-audio-streaming.git
cd omi-audio-streaming
go mod download
go run cmd/api/main.go


## Usage

1. Start the server:
   
   go run cmd/api/main.go
   

2. Send a POST request to `/process-audio` with a WAV file in the request body.

3. The server will process the audio, recognize the song, store it in Cloudinary, and save the track information in the database. -->

## API Endpoints

- `POST /audio`: Process audio bytes, recognize the song, and store it.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
