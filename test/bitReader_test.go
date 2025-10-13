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

func TestBitReader_Empty(t *testing.T) {
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
