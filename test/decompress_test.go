package test

import (
	"bytes"
	"huffman-compressor/internal"
	"os"
	"strings"
	"testing"
)

func TestDecompress_SingleCharacter(t *testing.T) {
	// Cretae test input file
	originalData := []byte("AAAA")
	originalPath := "test_roundtrip_original.txt"

	err := os.WriteFile(originalPath, originalData, 0644)
	if err != nil {
		t.Fatal("Failed to create original file:", err)
	}

	defer os.Remove(originalPath)

	// compress the file
	compressedPath := "test_roundtrip_compressed.hf"
	err = internal.CompressFile(originalPath, compressedPath)

	if err != nil {
		t.Fatal("Compression failed:", err)
	}

	defer os.Remove(compressedPath)

	// Decompress the file
	decompressedPath := "test_roundtrip_decompressed.txt"
	err = internal.Decompress(compressedPath, decompressedPath)
	if err != nil {
		t.Fatal("Decompression failed:", err)
	}

	defer os.Remove(decompressedPath)

	// Verify decompressed matched original
	err = internal.VerifyDecompression(originalPath, decompressedPath)
	if err != nil {
		t.Fatal("Verification failed:", err)
	}

	// Also manually verify content
	decompressedData, err := os.ReadFile(decompressedPath)
	if err != nil {
		t.Fatal("Failed to read decompressed file:", err)
	}

	if string(decompressedData) != string(originalData) {
		t.Errorf("Content mismatch: got %s, want %s",
			string(decompressedData), string(originalData))
	}

}

func TestDecompressFile_LargerTextRoundTrip(t *testing.T) {
	// Create larger file with good compression potential
	originalData := []byte(strings.Repeat("Hello, World! ", 100))
	originalPath := "test_roundtrip_large.txt"

	err := os.WriteFile(originalPath, originalData, 0644)
	if err != nil {
		t.Fatal("Failed to create original file:", err)
	}

	defer os.Remove(originalPath)

	compressedPath := "test_roundtrip_large_compressed.hf"
	err = internal.CompressFile(originalPath, compressedPath)
	if err != nil {
		t.Fatal("Compression failed:", err)
	}

	defer os.Remove(compressedPath)

	decompressedPath := "test_roundtrip_large_decompressed.txt"
	err = internal.Decompress(compressedPath, decompressedPath)
	if err != nil {
		t.Fatal("Decompression failed:", err)
	}

	defer os.Remove(decompressedPath)

	// Verify decompressed matched original
	err = internal.VerifyDecompression(originalPath, decompressedPath)
	if err != nil {
		t.Fatal("Verification failed:", err)
	}

	// log compression stats for info
	stats, err := internal.GetCompressionStats(originalPath, compressedPath)
	if err != nil {
		t.Fatal("Failed to get compression stats:", err)
	}

	internal.PrintCompressionStats(stats)

}

func TestDecompressFile_AllASCIIRoundTrip(t *testing.T) {
	// Cretae file with many diff characters
	originalData := []byte("The quick brown fox jumps over the lazy dog 0123456789!@#$%^&*()_+-=[]{}|;':\",./<>?")
	originalPath := "test_roundtrip_ascii.txt"

	err := os.WriteFile(originalPath, originalData, 0644)
	if err != nil {
		t.Fatal("Failed to create original file:", err)
	}

	defer os.Remove(originalPath)

	compressedPath := "test_roundtrip_ascii_compressed.hf"
	err = internal.CompressFile(originalPath, compressedPath)
	if err != nil {
		t.Fatal("Compression failed:", err)
	}

	defer os.Remove(compressedPath)

	decompressedPath := "test_roundtrip_ascii_decompressed.txt"
	err = internal.Decompress(compressedPath, decompressedPath)
	if err != nil {
		t.Fatal("Decompression failed:", err)
	}

	defer os.Remove(decompressedPath)

	// Verify decompressed matched original
	err = internal.VerifyDecompression(originalPath, decompressedPath)
	if err != nil {
		t.Fatal("Verification failed:", err)
	}
}

func TestDecompressFile_BinaryRoundTrip(t *testing.T) {
	//  Create binary file with various byte values
	originalData := make([]byte, 100)
	for i := 0; i < 100; i++ {
		originalData[i] = byte(i * 256)
	}

	originalPath := "test_roundtrip_binary.bin"
	err := os.WriteFile(originalPath, originalData, 0644)
	if err != nil {
		t.Fatal("Failed to create original file:", err)
	}

	defer os.Remove(originalPath)
	compressedPath := "test_roundtrip_binary_compressed.hf"
	err = internal.CompressFile(originalPath, compressedPath)
	if err != nil {
		t.Fatal("Compression failed:", err)
	}

	defer os.Remove(compressedPath)

	decompressedPath := "test_roundtrip_binary_decompressed.bin"
	err = internal.Decompress(compressedPath, decompressedPath)
	if err != nil {
		t.Fatal("Decompression failed:", err)
	}

	defer os.Remove(decompressedPath)
	// Verify decompressed matched original
	err = internal.VerifyDecompression(originalPath, decompressedPath)
	if err != nil {
		t.Fatal("Verification failed:", err)
	}
}

// ============================================================
// ERROR HANDLING TESTS
// ============================================================

func TestDecompressFile_NonEistentInput(t *testing.T) {
	err := internal.Decompress("non_existent_file.hf", "output.txt")

	if err == nil {
		t.Fatal("Expected error for non-existent input file, got nil")
		os.Remove("output.txt")
	}
}

func TestDecompressFile_InvalidMagicNumber(t *testing.T) {
	// create file with invalid magic number
	badPath := "test_bad_magic.hf"
	badData := []byte("XX") // Wrong magic, should be "HF"

	// Add some dummy data after magic
	badData = append(badData, make([]byte, 20)...)

	err := os.WriteFile(badPath, badData, 0644)
	if err != nil {
		t.Fatal("Failed to create bad magic file:", err)
	}

	defer os.Remove(badPath)

	outputPath := "test_bad_output.txt"
	err = internal.Decompress(badPath, outputPath)

	if err == nil {
		t.Fatal("Expected error for invalid magic number, got nil")
		os.Remove(outputPath)
	}

	if !strings.Contains(err.Error(), "magic") {
		t.Errorf("Unexpected error message about magic number, got: %v", err)
	}
}

func TestDecompresFile_CorruptData(t *testing.T) {
	// Create valid compressed file
	originalData := []byte("test dat for corruption")
	originalPath := "test_corrupt_original.txt"

	err := os.WriteFile(originalPath, originalData, 0644)
	if err != nil {
		t.Fatal("Failed to create original file:", err)
	}

	defer os.Remove(originalPath)

	compressedPath := "test_corrupt_compressed.hf"
	err = internal.CompressFile(originalPath, compressedPath)

	if err != nil {
		t.Fatal("Compression failed:", err)
	}

	defer os.Remove(compressedPath)

	// Corrupt the compressed file
	compressedData, err := os.ReadFile(compressedPath)
	if err != nil {
		t.Fatal("Failed to read compressed file:", err)
	}

	// Flip some bits in the middle of the file (after header)
	if len(compressedData) > 50 {
		compressedData[50] ^= 0xFF // bitwise OR ( alternate bits return 1 same bits return 0)
	}

	err = os.WriteFile(compressedPath, compressedData, 0644)
	if err != nil {
		t.Fatal("Failed to write corrupted file:", err)
	}

	// Try to decompress corrupted file
	decompressPath := "test_corrupt_output.txt"
	err = internal.Decompress(compressedPath, decompressPath)

	// should either error or produce incorrect output
	if err == nil {
		// if no error, verify it doesn't match original
		defer os.Remove(decompressPath)

		verifyErr := internal.VerifyDecompression(originalPath, decompressPath)
		if verifyErr == nil {
			t.Error("Corrupted file decompressed to correct output (unlikely!)")
		}
	}
	// if error occured, that's expected behavior
	t.Logf("Corrupted file properly rejected or produced wrong output")
}

// ============================================================
// BITREADER UNIT TESTS
// ============================================================

func TestBitReader_ReadBits(t *testing.T) {
	// create test data: byte 0b10110011
	testData := []byte{0b10110011}

	// create reader from bytes.Buffer or similar
	reader := bytes.NewReader(testData)
	bitReader := internal.NewBitReader(reader)

	// Expected bits (MSB first):: 1, 0, 1, 1, 0, 0, 1, 1
	expectedBits := []uint8{1, 0, 1, 1, 0, 0, 1, 1}

	for i := 0; i < 7; i++ {
		bit, err := bitReader.ReadBit()
		if err != nil {
			t.Fatalf("Failed to read bit %d: %v", i, err)
		}
		if bit != expectedBits[i] {
			t.Errorf("Bit %d: expected %d, got %d", i, expectedBits[i], bit)
		}
	}
}

func TestBitReader_MultipleBytes(t *testing.T) {
	// create test data: two bytes
	testData := []byte{0b11110000, 0b10101010}

	reader := bytes.NewReader(testData)
	bitReader := internal.NewBitReader(reader)

	// Expectedd: 1111 0000 1010 1010 (16 bits)
	expectedBits := []uint8{
		1, 1, 1, 1, 0, 0, 0, 0, // First byte
		1, 0, 1, 0, 1, 0, 1, 0, // Second byte
	}

	for i := 0; i < 16; i++ {
		bit, err := bitReader.ReadBit()
		if err != nil {
			t.Fatalf("Failed to read bit %d: %v", i, err)
		}

		if bit != expectedBits[i] {
			t.Errorf("Bit %d: expected %d, got %d", i, expectedBits[i], bit)
		}
	}
}

// ============================================================
// VERIFICATION FUNCTION TESTS
// ============================================================

func TestVerifyDecompression_Identical(t *testing.T) {
	// Create two identical files
	data := []byte("identical content")

	file1 := "test_verify_1.txt"
	file2 := "test_verify_2.txt"

	os.WriteFile(file1, data, 0644)
	os.WriteFile(file2, data, 0644)

	defer os.Remove(file1)
	defer os.Remove(file2)

	err := internal.VerifyDecompression(file1, file2)
	if err != nil {
		t.Errorf("Verification failed for identical files: %v", err)
	}
}

func TestVerifyDecompression_Different(t *testing.T) {
	// create two different files
	file1 := "test_verify_diff1.txt"
	file2 := "test_verify_diff2.txt"

	os.WriteFile(file1, []byte("content A"), 0644)
	os.WriteFile(file2, []byte("content B"), 0644)

	defer os.Remove(file1)
	defer os.Remove(file2)

	err := internal.VerifyDecompression(file1, file2)
	if err == nil {
		t.Error("Expected error for different files, got nil")
	}
}

func TestVerifyDecompression_DifferentSizes(t *testing.T) {
	file1 := "test_verify_size1.txt"
	file2 := "test_verify_size2.txt"

	os.WriteFile(file1, []byte("short"), 0644)
	os.WriteFile(file2, []byte("much longer content"), 0644)

	defer os.Remove(file1)
	defer os.Remove(file2)

	err := internal.VerifyDecompression(file1, file2)
	if err == nil {
		t.Error("Expected error for different sizes, got nil")
	}

	if !strings.Contains(err.Error(), "size mismatch") {
		t.Errorf("Expected size mismatch error, got: %v", err)
	}
}

// ============================================================
// EDGE CASE TESTS
// ============================================================
func TestDecompressFile_SpecialChars(t *testing.T) {
	originalData := []byte("Line 1\nLine 2 \r\nLine 3\tTabbedx00Null")
	originalPath := "test_special.txt"

	os.WriteFile(originalPath, originalData, 0644)
	defer os.Remove(originalPath)

	compressedPath := "test_special.hf"
	internal.CompressFile(originalPath, compressedPath)
	defer os.Remove(compressedPath)

	decompressedPath := "test_special_dec.txt"
	internal.Decompress(compressedPath, decompressedPath)
	defer os.Remove(decompressedPath)

	err := internal.VerifyDecompression(originalPath, decompressedPath)
	if err != nil {
		t.Fatal("Special characters not preserved:", err)
	}
}
