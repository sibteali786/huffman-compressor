package internal

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

const MagicNumber = "HF"

type FileHeader struct {
	OriginalSize uint64         // Original uncompressed file size
	NumChars     uint8          // Number of unique characters
	PaddingBits  uint8          // Number of padding bits in last byte
	FreqTable    FrequencyTable // Character frequencies
}

func WriteHeader(writer io.Writer, freqTable FrequencyTable, originalSize uint64, paddingBits uint8) error {
	// 1. Write magic number "HF"
	_, err := writer.Write([]byte(MagicNumber))
	if err != nil {
		return err
	}

	// 2. Write original file size (8 bytes, big-endian)
	err = binary.Write(writer, binary.BigEndian, originalSize)
	if err != nil {
		return err
	}

	// 3. Write number of unique characters
	numChars := uint8(len(freqTable))
	_, err = writer.Write([]byte{numChars})
	if err != nil {
		return err
	}

	// 4. Write padding bits count
	_, err = writer.Write([]byte{paddingBits})
	if err != nil {
		return err
	}

	// ✅ 5. Write frequency entries IN SORTED ORDER
	// Sort characters to ensure deterministic order
	chars := make([]byte, 0, len(freqTable))
	for char := range freqTable {
		chars = append(chars, char)
	}
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})

	// Write in sorted order
	for _, char := range chars {
		freq := freqTable[char]
		_, err = writer.Write([]byte{char})
		if err != nil {
			return err
		}

		err = binary.Write(writer, binary.BigEndian, uint32(freq))
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadHeader reads and parses the header from reader
func ReadHeader(reader io.Reader) (FileHeader, error) {
	header := FileHeader{}
	// 1. Read and verify magic number
	magic := make([]byte, 2)
	_, err := io.ReadFull(reader, magic)
	if err != nil {
		return header, err
	}

	if string(magic) != MagicNumber {
		return header, fmt.Errorf("invalid file format: bad magic number")
	}

	// 2. Read original size
	err = binary.Read(reader, binary.BigEndian, &header.OriginalSize)
	if err != nil {
		return header, err
	}

	// 3. Read number of unique characters
	numChars, err := readByte(reader)
	if err != nil {
		return header, err
	}
	header.NumChars = numChars

	// 4. Read padding bits count
	paddingBits, err := readByte(reader)
	if err != nil {
		return header, err
	}
	header.PaddingBits = paddingBits

	// 5. Read frequency entries
	header.FreqTable = make(FrequencyTable)
	for i := 0; i < int(header.NumChars); i++ {
		char, err := readByte(reader)
		if err != nil {
			return header, err
		}

		var freq uint32
		err = binary.Read(reader, binary.BigEndian, &freq)
		if err != nil {
			return header, err
		}

		header.FreqTable[char] = int(freq)
	}

	return header, nil
}

func readByte(reader io.Reader) (byte, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func CalculateHeaderSize(freqTable FrequencyTable) int {
	// Magic(2) + OriginalSize(8) + NumChars(1) + PaddingBits(1) + Entries(N*5)
	return 12 + (len(freqTable) * 5)
}
