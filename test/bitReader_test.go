package test

import (
	"bytes"
	"huffman-compressor/internal"
	"io"
	"testing"
)

func TestBitReader_ReadSingleByte(t *testing.T) {
	// Prepare a byte buffer with known data
	data := []byte{0b10110011} // 8 bits: 10110011
	reader := bytes.NewReader(data)
	br := internal.NewBitReader(reader)

	// Read bits one by one and verify
	expectedBits := []uint8{1, 0, 1, 1, 0, 0, 1, 1}
	for i, expected := range expectedBits {
		bit, err := br.ReadBit()
		if err != nil {
			t.Fatal("Unexpected error:", err)
		}
		if bit != expected {
			t.Errorf("Bit %d: expected %d, got %d", i, expected, bit)
		}
	}

	// Next read should return EOF
	_, err := br.ReadBit()
	if err != io.EOF {
		t.Error("Expected EOF, got:", err)
	}

	if !br.IsFinished() {
		t.Error("Expected IsFinished to be true after EOF")
	}

}

func TestBitReader_ReadMultipleBytes(t *testing.T) {
	// Prepare a byte buffer with known data
	data := []byte{0b11001100, 0b10101010} // 16 bits: 11001100 10101010
	reader := bytes.NewReader(data)
	br := internal.NewBitReader(reader)

	// Read bits one by one and verify
	expectedBits := []uint8{
		1, 1, 0, 0, 1, 1, 0, 0, // First byte
		1, 0, 1, 0, 1, 0, 1, 0, // Second byte
	}
	for i, expected := range expectedBits {
		bit, err := br.ReadBit()
		if err != nil {
			t.Fatal("Unexpected error:", err)
		}
		if bit != expected {
			t.Errorf("Bit %d: expected %d, got %d", i, expected, bit)
		}
	}

	// Next read should
	_, err := br.ReadBit()
	if err != io.EOF {
		t.Error("Expected EOF, got:", err)
	}

	if !br.IsFinished() {
		t.Error("Expected IsFinished to be true after EOF")
	}

}

func TestBitReader_EmptyInput(t *testing.T) {
	// Prepare an empty byte buffer
	data := []byte{}
	reader := bytes.NewReader(data)
	br := internal.NewBitReader(reader)

	// Attempt to read a bit, should get EOF immediately
	_, err := br.ReadBit()
	if err != io.EOF {
		t.Error("Expected EOF, got:", err)
	}
	if !br.IsFinished() {
		t.Error("Expected IsFinished to be true after EOF")
	}
}

// TestBitReader_AllZeros tests reading a byte of all zeros
func TestBitReader_AllZeros(t *testing.T) {
	{
		data := []byte{0b00000000}
		reader := bytes.NewReader(data)
		bitReader := internal.NewBitReader(reader)

		for i := 0; i < 8; i++ {
			bit, err := bitReader.ReadBit()
			if err != nil {
				t.Fatalf("Unexpected error at bit %d: %v", i, err)
			}
			if bit != 0 {
				t.Errorf("Bit %d: expected 0, got %d", i, bit)
			}
		}
	}
}

// TestBitReader_AllOnes tests reading a byte of all ones
func TestBitReader_AllOnes(t *testing.T) {
	data := []byte{0b11111111}
	reader := bytes.NewReader(data)
	bitReader := internal.NewBitReader(reader)

	for i := 0; i < 8; i++ {
		bit, err := bitReader.ReadBit()
		if err != nil {
			t.Fatalf("Unexpected error at bit %d: %v", i, err)
		}
		if bit != 1 {
			t.Errorf("Bit %d: expected 1, got %d", i, bit)
		}
	}
}

// TestBitReader_SpecificPattern tests a specific bit pattern
func TestBitReader_SpecificPattern(t *testing.T) {
	/// Byte: 0b10101010 = 170 (alternating bits)
	data := []byte{170}
	reader := bytes.NewReader(data)
	bitReader := internal.NewBitReader(reader)

	expected := []uint8{1, 0, 1, 0, 1, 0, 1, 0}

	for i, expectedBit := range expected {
		bit, err := bitReader.ReadBit()
		if err != nil {
			t.Fatalf("Unexpected error at bit %d: %v", i, err)
		}
		if bit != expectedBit {
			t.Errorf("Bit %d: expected %d, got %d", i, expectedBit, bit)
		}
	}

}

// TestBitReader_LargeInput tests reading from a larger byte stream
func TestBitReader_LargeInput(t *testing.T) {
	// Create 10 bytes of data
	data := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	reader := bytes.NewReader(data)
	bitReader := internal.NewBitReader(reader)

	bitCount := 0
	for {
		_, err := bitReader.ReadBit()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		bitCount++
	}

	expectedBits := len(data) * 8
	if bitCount != expectedBits {
		t.Errorf("Expected to read %d bits, got %d", expectedBits, bitCount)
	}
}

// TestBitReader_ReconstructByte tests reading 8 bits and reconstructing the original byte
func TestBitReader_ReconstructByte(t *testing.T) {
	testCases := []byte{
		0b00000000, // 0
		0b11111111, // 255
		0b10101010, // 170
		0b01010101, // 85
		0b11110000, // 240
		0b00001111, // 15
		0b10011001, // 153
	}

	for _, originalByte := range testCases {
		data := []byte{originalByte}
		reader := bytes.NewReader(data)
		bitReader := internal.NewBitReader(reader)

		// Read 8 bits and reconstruct the byte
		var reconstructed byte = 0
		for i := 0; i < 8; i++ {
			bit, err := bitReader.ReadBit()
			if err != nil {
				t.Fatalf("Error reading bit %d for byte %d: %v", i, originalByte, err)
			}
			reconstructed = (reconstructed << 1) | bit
		}

		if reconstructed != originalByte {
			t.Errorf("Byte reconstruction failed: original=%08b (%d), reconstructed=%08b (%d)",
				originalByte, originalByte, reconstructed, reconstructed)
		}
	}
}

// TestBitReader_PartialByteRead tests reading only some bits from a byte
func TestBitReader_PartialByteRead(t *testing.T) {
	data := []byte{0b11110000}
	reader := bytes.NewReader(data)
	bitReader := internal.NewBitReader(reader)

	// Read only first 4 bits
	expectedBits := []uint8{1, 1, 1, 1}
	for i, expected := range expectedBits {
		bit, err := bitReader.ReadBit()
		if err != nil {
			t.Fatalf("Unexpected error at bit %d: %v", i, err)
		}
		if bit != expected {
			t.Errorf("Bit %d: expected %d, got %d", i, expected, bit)
		}
	}

	// We should still be able to read the remaining 4 bits
	for i := 0; i < 4; i++ {
		bit, err := bitReader.ReadBit()
		if err != nil {
			t.Fatalf("Unexpected error reading remaining bits: %v", err)
		}
		if bit != 0 {
			t.Errorf("Expected remaining bits to be 0, got %d", bit)
		}
	}
}

// TestBitReader_IsFinishedBeforeEOF tests that IsFinished is false before EOF
func TestBitReader_IsFinishedBeforeEOF(t *testing.T) {
	data := []byte{0b10101010}
	reader := bytes.NewReader(data)
	bitReader := internal.NewBitReader(reader)

	// IsFinished should be false initially
	if bitReader.IsFinished() {
		t.Error("Expected IsFinished() to return false initially")
	}

	// Read a few bits
	for i := 0; i < 4; i++ {
		_, err := bitReader.ReadBit()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	// Should still not be finished
	if bitReader.IsFinished() {
		t.Error("Expected IsFinished() to return false after partial read")
	}

	// Read remaining bits
	for i := 0; i < 4; i++ {
		_, err := bitReader.ReadBit()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	// Try to read one more (should get EOF)
	_, err := bitReader.ReadBit()
	if err != io.EOF {
		t.Errorf("Expected EOF, got: %v", err)
	}

	// Now should be finished
	if !bitReader.IsFinished() {
		t.Error("Expected IsFinished() to return true after EOF")
	}
}

// TestBitReader_MultipleEOFReads tests that multiple reads after EOF continue to return EOF
func TestBitReader_MultipleEOFReads(t *testing.T) {
	data := []byte{0b11111111}
	reader := bytes.NewReader(data)
	bitReader := internal.NewBitReader(reader)

	// Read all 8 bits
	for i := 0; i < 8; i++ {
		bitReader.ReadBit()
	}

	// Multiple EOF reads should all return EOF
	for i := 0; i < 5; i++ {
		_, err := bitReader.ReadBit()
		if err != io.EOF {
			t.Errorf("Read %d after EOF: expected EOF, got %v", i+1, err)
		}
		if !bitReader.IsFinished() {
			t.Errorf("Read %d after EOF: IsFinished() should return true", i+1)
		}
	}
}
