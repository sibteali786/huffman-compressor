package test

import (
	"bytes"
	"huffman-compressor/internal"
	"testing"
)

func TestWriteAndReadHeader(t *testing.T) {
	freqTable := internal.FrequencyTable{
		'a': 100,
		'b': 50,
		'c': 25,
	}
	originalSize := uint64(175)
	paddingBits := uint8(3)

	// Write header to buffer
	var buffer bytes.Buffer
	err := internal.WriteHeader(&buffer, freqTable, originalSize, paddingBits)
	if err != nil {
		t.Fatal("Failed to write header:", err)
	}

	// Verify header size
	header, err := internal.ReadHeader(&buffer)
	if err != nil {
		t.Fatal("Failed to read header:", err)
	}

	// Verify all fields match
	if header.OriginalSize != originalSize {
		t.Errorf("OriginalSize mismatch: got %d, want %d", header.OriginalSize, originalSize)
	}
	if header.NumChars != uint8(len(freqTable)) {
		t.Errorf("NumChars mismatch: got %d, want %d", header.NumChars, len(freqTable))
	}
	if header.PaddingBits != paddingBits {
		t.Errorf("PaddingBits mismatch: got %d, want %d", header.PaddingBits, paddingBits)
	}
	if len(header.FreqTable) != len(freqTable) {
		t.Errorf("FreqTable length mismatch: got %d, want %d", len(header.FreqTable), len(freqTable))
	}

	if header.FreqTable['a'] != freqTable['a'] {
		t.Errorf("FreqTable['a'] mismatch: got %d, want %d", header.FreqTable['a'], freqTable['a'])
	}
	if header.FreqTable['b'] != freqTable['b'] {
		t.Errorf("FreqTable['b'] mismatch: got %d, want %d", header.FreqTable['b'], freqTable['b'])
	}
	if header.FreqTable['c'] != freqTable['c'] {
		t.Errorf("FreqTable['c'] mismatch: got %d, want %d", header.FreqTable['c'], freqTable['c'])
	}
}

func TestReadHeader_InvalidMagic(t *testing.T) {
	// Write invalid magic number to buffer
	var buffer bytes.Buffer
	buffer.Write([]byte("XX"))

	// Try to read
	_, err := internal.ReadHeader(&buffer)
	if err == nil {
		t.Fatal("Expected error for invalid magic number, got nil")
	}

	if err.Error() != "invalid file format: bad magic number" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestHeaderSize_Calculation(t *testing.T) {
	tests := []struct {
		numChars int
		expected int
	}{
		{1, 17},
		{10, 62},
		{256, 1292},
	}

	for _, test := range tests {
		freqTable := make(internal.FrequencyTable)
		for i := 0; i < test.numChars; i++ {
			freqTable[byte(i)] = 1
		}
		size := internal.CalculateHeaderSize(freqTable)
		if size != test.expected {
			t.Errorf("Header size calculation failed for %d chars: got %d, want %d", test.numChars, size, test.expected)
		}
	}
}

func TestEmptyFrequency(t *testing.T) {
	freqTable := internal.FrequencyTable{}

	var buffer bytes.Buffer
	err := internal.WriteHeader(&buffer, freqTable, 0, 0)
	if err != nil {
		t.Fatal("Failed to write header:", err)
	}

	if buffer.Len() != 12 {
		t.Errorf("Expected header size 12 for empty freq table, got %d", buffer.Len())
	}
}

func TestSingleCharacter(t *testing.T) {
	freqTable := internal.FrequencyTable{'x': 100}

	var buffer bytes.Buffer
	err := internal.WriteHeader(&buffer, freqTable, 100, 1)
	if err != nil {
		t.Fatal("Failed to write header:", err)
	}

	header, err := internal.ReadHeader(&buffer)
	if err != nil {
		t.Fatal("Failed to read header:", err)
	}

	if header.NumChars != 1 {
		t.Errorf("Expected NumChars 1, got %d", header.NumChars)
	}
	if header.FreqTable['x'] != 100 {
		t.Errorf("Expected frequency of 'x' to be 100, got %d", header.FreqTable['x'])
	}

}
