package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// PCMToWAV wraps raw 16-bit mono 16kHz PCM audio bytes in a standard 44-byte WAV header.
func PCMToWAV(pcm []byte, sampleRate int, channels int, bitsPerSample int) ([]byte, error) {
	if sampleRate <= 0 {
		sampleRate = 16000
	}
	if channels <= 0 {
		channels = 1
	}
	if bitsPerSample <= 0 {
		bitsPerSample = 16
	}

	dataSize := uint32(len(pcm))
	fileSize := 36 + dataSize
	byteRate := uint32(sampleRate * channels * (bitsPerSample / 8))
	blockAlign := uint16(channels * (bitsPerSample / 8))

	buf := new(bytes.Buffer)

	// RIFF header
	buf.WriteString("RIFF")
	if err := binary.Write(buf, binary.LittleEndian, fileSize); err != nil {
		return nil, err
	}
	buf.WriteString("WAVE")

	// fmt subchunk
	buf.WriteString("fmt ")
	if err := binary.Write(buf, binary.LittleEndian, uint32(16)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint16(1)); err != nil { // 1 = PCM
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint16(channels)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint32(sampleRate)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, byteRate); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, blockAlign); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint16(bitsPerSample)); err != nil {
		return nil, err
	}

	// data subchunk
	buf.WriteString("data")
	if err := binary.Write(buf, binary.LittleEndian, dataSize); err != nil {
		return nil, err
	}

	// append PCM audio payload
	if _, err := buf.Write(pcm); err != nil {
		return nil, fmt.Errorf("failed to write PCM data: %w", err)
	}

	return buf.Bytes(), nil
}
