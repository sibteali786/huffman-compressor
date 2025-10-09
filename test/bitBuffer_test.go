package test

import (
	"bytes"
	"huffman-compressor/internal"
	"testing"
)

func TestWriteSingleBit(t *testing.T) {
	// Create a buffer to capture output
	var output bytes.Buffer

	bb := internal.NewBitBuffer(&output)

	// Write a single bit '1'
	bb.WriteBit(1)

	// Close to flush
	bb.Close()

	// Verify
	result := output.Bytes()
	expected := []byte{0b10000000} // "1" followed by 7 padding zeros

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestWritePartialBits(t *testing.T) {
	var output bytes.Buffer
	bb := internal.NewBitBuffer(&output)
	// Write "101" ( 3 bits )
	bb.WriteBit(1)
	bb.WriteBit(0)
	bb.WriteBit(1)

	bb.Close()

	result := output.Bytes()
	expected := []byte{0b10100000} // "101" + 5 padding zeros

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %08b, got %08b", expected[0], result[0])
	}
}

// Test: Writing exactly 8 bits (one complete byte)
func TestWriteCompleteByte(t *testing.T) {

	var output bytes.Buffer
	bb := internal.NewBitBuffer(&output)

	// Write "10110011" (8 bits)
	bits := []int{1, 0, 1, 1, 0, 0, 1, 1}
	for i := range bits {
		bb.WriteBit(uint64(bits[i]))
	}

	bb.Close()

	result := output.Bytes()
	expected := []byte{0b10110011} // Complete byte, no padding

	if len(result) != 1 {

		t.Errorf("Expected 1 byte, got %d", len(result))
	}

	if result[0] != expected[0] {

		t.Errorf("Expected %08b, got %08b", expected[0], result[0])
	}
}

// Test: Writing more than 8 bits (multiple bytes)
func TestWriteMultipleBytes(t *testing.T) {
	var output bytes.Buffer
	bb := internal.NewBitBuffer(&output)

	// Write "10110011" + "11001" (13 bits = 1 full byte + 5 bits)
	bits := []int{1, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1}
	for i := range bits {
		bb.WriteBit(uint64(bits[i]))
	}

	bb.Close()

	result := output.Bytes()
	expected := []byte{
		0b10110011, // First 8 bits
		0b11001000, // Last 5 bits + 3 padding zeros
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 bytes, got %d", len(result))
	}

	for i := range len(expected) - 1 {
		if result[i] != expected[i] {
			t.Errorf("Byte %d: expected %08b, got %08b", i, expected[i], result[i])
		}
	}
}

// Test: WriteBits function with HuffmanCode-like structure
func TestWriteBitsFunction(t *testing.T) {
	var output bytes.Buffer
	bb := internal.NewBitBuffer(&output)
	// Write code "110" (3 bits, value 6) in MSB Format
	// But we write in In LSB-first storage: bits=3, length=3
	bb.WriteBits(3, 3)

	bb.Close()

	result := output.Bytes()
	expected := []byte{0b11000000} // "110" + 5 padding zeros

	if result[0] != expected[0] {
		t.Errorf("Expected %08b, got %08b", expected[0], result[0])
	}
}

// Test: Writing multiple codes (simulating Huffman encoding)
func TestWriteMultipleCodes(t *testing.T) {

	var output bytes.Buffer
	bb := internal.NewBitBuffer(&output)

	// Simulate encoding "ABC" with codes:
	// 'A' = "10"   (2 bits, value 2) => 01
	// 'B' = "110"  (3 bits, value 6) => 011
	// 'C' = "1111" (4 bits, value 15) =>

	bb.WriteBits(1, 2)  // "01"
	bb.WriteBits(3, 3)  // "011"
	bb.WriteBits(15, 4) // "1111"

	// Total: 2+3+4 = 9 bits
	// Should produce: "10" + "110" + "1111" = "101101111"
	// = "10110111" + "1" + 7 padding zeros

	bb.Close()

	result := output.Bytes()
	expected := []byte{
		0b10110111, // First 8 bits
		0b10000000, // Last 1 bit + 7 padding zeros
	}

	if len(result) != 2 {

		t.Errorf("Expected 2 bytes, got %d", len(result))
	}

	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] {
			t.Errorf("Byte %d: expected %08b, got %08b", i, expected[i], result[i])
		}

	}
}

func TestWriteOnlyZeros(t *testing.T) {
	var output bytes.Buffer
	bb := internal.NewBitBuffer(&output)

	// Write "000" (3 zeros)
	bb.WriteBit(0)
	bb.WriteBit(0)
	bb.WriteBit(0)

	bb.Close()

	result := output.Bytes()
	expected := []byte{0b00000000} // All zeros

	if result[0] != expected[0] {
		t.Errorf("Expected %08b, got %08b", expected[0], result[0])
	}
}

// Test: Alternating bits pattern
func TestAlternatingBits(t *testing.T) {
	var output bytes.Buffer

	bb := internal.NewBitBuffer(&output)

	// Write "10101010" (8 bits alternating)
	for i := 0; i < 8; i++ {
		bb.WriteBit(uint64((i + 1) % 2))
	}

	bb.Close()

	result := output.Bytes()
	expected := []byte{0b10101010}

	if result[0] != expected[0] {
		t.Errorf("Expected %08b, got %08b", expected[0], result[0])
	}
}

// Test: Total bits counter
func TestTotalBitsCounter(t *testing.T) {

	var output bytes.Buffer
	bb := internal.NewBitBuffer(&output)

	// Write 13 bits
	bb.WriteBits(0xFF, 8) // 8 bits
	bb.WriteBits(0x1F, 5) // 5 bits

	totalBits := bb.GetTotalBits()

	if totalBits != 13 {

		t.Errorf("Expected 13 total bits, got %d", totalBits)
	}

	bb.Close()
}

// Test: large data ( stress test )
func TestLargeData(t *testing.T) {
	var output bytes.Buffer

	bb := internal.NewBitBuffer(&output)
	// Write 1000 bits
	for i := 0; i < 1000; i++ {
		bb.WriteBit(uint64((i + 1) % 2))
	}

	bb.Close()

	result := output.Bytes()
	expectedBytes := 1000 / 8 // 125 bytes

	if len(result) != expectedBytes {
		t.Errorf("Expected %d bytes, got %d", expectedBytes, len(result))
	}
	// Verify pattern (alternating 0b10101010)
	for i := 0; i < len(result); i++ {
		if result[i] != 0b10101010 {
			t.Errorf("Byte %d: expected %08b, got %08b", i, 0b10101010, result[i])
		}
	}
}

// Test: Flush behavior (byte written immediately)
func TestFlushBehavior(t *testing.T) {

	var output bytes.Buffer
	bb := internal.NewBitBuffer(&output)

	// Write 7 bits - should NOT flush yet
	for i := 0; i < 7; i++ {
		bb.WriteBit(1)
	}

	if len(output.Bytes()) != 0 {
		t.Errorf("Should not flush incomplete byte")
	}

	// Write 8th bit - should flush
	bb.WriteBit(1)

	if len(output.Bytes()) != 1 {

		t.Errorf("Should flush after 8 bits")
	}

	result := output.Bytes()
	expected := byte(0b11111111)

	if result[0] != expected {
		t.Errorf("Expected %08b, got %08b", expected, result[0])
	}

}

// Test: Writing with different bit lengths
func TestVariableLengthCodes(t *testing.T) {
	testCases := []struct {
		bits     uint64
		length   int
		expected string
	}{
		{1, 1, "1"},     // "1" → LSB: pos0=1 → bits=1 ✓
		{1, 2, "10"},    // "10" → LSB: pos0=1, pos1=0 → bits=1 (not 2!)
		{7, 3, "111"},   // "111" → LSB: pos0=1, pos1=1, pos2=1 → bits=7 ✓
		{15, 4, "1111"}, // "1111" → LSB: all 1s → bits=15 ✓
		{0, 3, "000"},   // "000" → LSB: all 0s → bits=0 ✓
	}

	for _, tc := range testCases {
		var output bytes.Buffer
		bb := internal.NewBitBuffer(&output)

		bb.WriteBits(tc.bits, tc.length)
		bb.Close()

		result := output.Bytes()

		if !CompareBits(result[0], tc.expected, tc.length) {
			t.Errorf("Code %s: expected byte %08b, got %08b",
				tc.expected, stringToByte(tc.expected), result[0])
		}

		t.Logf("Code: %s wrote byte: %08b", tc.expected, result[0])
	}
}

func stringToByte(s string) any {
	var b byte = 0
	for i := 0; i < len(s); i++ {
		bitPos := 7 - i
		if s[i] == '1' {
			b |= (1 << bitPos)
		}
	}
	return b
}

// Helper function to compare bits
func CompareBits(actual byte, expectedBits string, length int) bool {

	for i := 0; i < length; i++ {
		bitPos := 7 - i
		actualBit := (actual >> bitPos) & 1
		expectedBit := 0
		if expectedBits[i] == '1' {
			expectedBit = 1
		}

		if int(actualBit) != expectedBit {

			return false
		}

	}
	return true
}

func TestBitOrderingExplicit(t *testing.T) {
	var output bytes.Buffer
	bb := internal.NewBitBuffer(&output)

	// bits = 0b0101 = 5, length = 4
	// Should write: 0, 1, 0, 1
	bb.WriteBits(5, 4)
	bb.Close()

	result := output.Bytes()
	expected := byte(0b01010000) // "0101" + padding

	if result[0] != expected {
		t.Errorf("Expected %08b, got %08b", expected, result[0])
	}
}
