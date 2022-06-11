package wav

import (
	"bytes"
	"encoding/binary"
)

type Wav struct {
	// "RIFF" chunk descriptor
	ChunkID   uint32 // "RIFF"
	ChunkSize uint32 // 36 + SubChunk2Size
	Format    uint32 // "WAVE"

	// "fmt" sub-chunk
	Subchunk1ID   uint32 // "fmt "
	Subchunk1Size uint32 // sum of the rest subchunk size (2+2+4+4+2+2=16)
	AudioFormat   uint16 // 1 for PCM
	NumChannels   uint16 // 2 for stereo
	SampleRate    uint32 // 48000 Hz (80 BB 00 00)
	ByteRate      uint32 // ByteRates=(Sample Rate x Bits Per Sample x Channel Numbers)/8
	BlockAlign    uint16 // Data block size
	BitsPerSample uint16 // 16bits

	// "data" sub-chunk
	Subchunk2ID   uint32 // "data"
	Subchunk2Size uint32 // Number of bytes in the data (Sample numbers x Channel numbers x Bits per sample)/8

	data []byte
}

func SetWav(SampleRate uint32, BitsPerSample uint16, data []byte) Wav {
	var w Wav

	w.ChunkID = uint32(0x52494646) // "RIFF"
	w.ChunkSize = uint32(36 + len(data))
	w.Format = uint32(0x57415645) // "WAVE"

	w.Subchunk1ID = uint32(0x666d7420) // "fmt "
	w.Subchunk1Size = uint32(16)
	w.AudioFormat = uint16(1)
	w.NumChannels = uint16(2)
	w.SampleRate = SampleRate
	w.ByteRate = (SampleRate * uint32(BitsPerSample) * uint32(w.NumChannels)) / 8
	w.BlockAlign = (w.NumChannels * BitsPerSample) / 8
	w.BitsPerSample = BitsPerSample

	w.Subchunk2ID = uint32(0x64617461) // "data"
	w.Subchunk2Size = uint32(len(data))
	w.data = data

	return w
}

func (w Wav) Encode() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, w.ChunkID)
	binary.Write(buf, binary.LittleEndian, w.ChunkSize)
	binary.Write(buf, binary.BigEndian, w.Format)

	binary.Write(buf, binary.BigEndian, w.Subchunk1ID)
	binary.Write(buf, binary.LittleEndian, w.Subchunk1Size)
	binary.Write(buf, binary.LittleEndian, w.AudioFormat)
	binary.Write(buf, binary.LittleEndian, w.NumChannels)
	binary.Write(buf, binary.LittleEndian, w.SampleRate)
	binary.Write(buf, binary.LittleEndian, w.ByteRate)
	binary.Write(buf, binary.LittleEndian, w.BlockAlign)
	binary.Write(buf, binary.LittleEndian, w.BitsPerSample)

	binary.Write(buf, binary.BigEndian, w.Subchunk2ID)
	binary.Write(buf, binary.LittleEndian, w.Subchunk2Size)

	binary.Write(buf, binary.LittleEndian, w.data)

	return buf.Bytes()
}
