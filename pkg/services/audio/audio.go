package audio

import (
    // "bytes"
    "encoding/binary"
    "fmt"
    "math"
)

const (
    SilenceThreshold = 1000 // Adjust based on background noise level
    SampleRate       = 16000 // Common sample rate; adjust if needed
    BitsPerSample    = 16    // For 16-bit PCM
    Channels         = 1     // Mono audio
)

// CheckIfContainsSoundRaw checks raw PCM audio data for sound by analyzing sample amplitude
func CheckIfContainsSoundRaw(data []byte) (bool, error) {
    sampleSize := BitsPerSample / 8 * Channels // Size of each sample in bytes
    if len(data)%sampleSize != 0 {
        return false, fmt.Errorf("audio data length is not a multiple of sample size")
    }

    totalSamples := len(data) / sampleSize
    significantSamples := 0
    threshold := SilenceThreshold

    // Iterate over samples
    for i := 0; i < totalSamples; i++ {
        // Read 16-bit sample (assuming signed little-endian PCM)
        sample := int16(binary.LittleEndian.Uint16(data[i*sampleSize : i*sampleSize+2]))

        // Check if the sample amplitude is above the silence threshold
        if math.Abs(float64(sample)) > float64(threshold) {
            significantSamples++
        }

        // If a significant number of samples contain sound, return true
        if significantSamples > totalSamples/20 { // Adjust this ratio as needed
            return true, nil
        }
    }
    return false, nil
}

// func main() {
//     // Example raw PCM data (replace this with actual data)
//     rawAudio := []byte{ /* Your raw PCM audio bytes here */ }

//     containsSound, err := CheckIfContainsSoundRaw(rawAudio)
//     if err != nil {
//         fmt.Println("Error:", err)
//     } else if containsSound {
//         fmt.Println("The audio contains sound!")
//     } else {
//         fmt.Println("The audio is silent.")
//     }
// }