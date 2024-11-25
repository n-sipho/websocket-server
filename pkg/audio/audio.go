package audio

import "encoding/binary"

const (
	numChannels   = 1 // Mono audio
	sampleRate    = 16000
	bitsPerSample = 16 // 16 bits per sample
)

// CreateWAVHeader generates a WAV header for the given data length
func CreateWAVHeader(dataLength int) []byte {
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