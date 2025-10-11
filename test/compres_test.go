package test

import (
	"huffman-compressor/internal"
	"os"
	"strings"
	"testing"
)

// Test: Compress simple known text
func TestCompressFile_SimpleText(t *testing.T) {
	// Create test input file
	testData := []byte("aaabbc")
	inputPath := "test_compress_input.txt"

	err := os.WriteFile(inputPath, testData, 0644)
	if err != nil {
		t.Fatal("Failed to create test file: ", err)
	}

	defer os.Remove(inputPath)

	// Compress
	outputPath := "test_compress_output.hf"
	err = internal.CompressFile(inputPath, outputPath)
	if err != nil {
		t.Fatal("compression failed: ", err)
	}

	defer os.Remove(outputPath)

	// verify the output file exists
	outputFile, err := os.Open(outputPath)
	if err != nil {
		t.Fatal("Failed to open output file: ", err)
	}
	defer outputFile.Close()

	header, err := internal.ReadHeader(outputFile)
	if err != nil {
		t.Fatal("Failed to read header: ", err)
	}

	// Verify header contents
	if header.OriginalSize != 6 {
		t.Errorf("Expected original size %d, got %d", 6, header.OriginalSize)
	}

	if len(header.FreqTable) != 3 {
		t.Errorf("Expected %d unique chars, got %d", 3, len(header.FreqTable))
	}

	if header.FreqTable['a'] != 3 {
		t.Errorf("Expected frequency of 'a' to be %d, got %d", 3, header.FreqTable['a'])
	}
	if header.FreqTable['b'] != 2 {
		t.Errorf("Expected frequency of 'b' to be 2, got %d", header.FreqTable['b'])
	}
	if header.FreqTable['c'] != 1 {
		t.Errorf("Expected frequency of 'c' to be 1, got %d", header.FreqTable['c'])
	}

	// Verify padding bits is set (should be 0-7)
	if header.PaddingBits > 7 {
		t.Errorf("Invalid padding bits: %d", header.PaddingBits)
	}
}

func TestCompressFile_EmptyFile(t *testing.T) {
	inputPath := "test_compress_empty.txt"

	err := os.WriteFile(inputPath, []byte{}, 0644)
	if err != nil {
		t.Fatal("Failed to create empty file: ", err)
	}
	defer os.Remove(inputPath)

	outputPath := "test_compress_empty.hf"

	err = internal.CompressFile(inputPath, outputPath)
	if err == nil {
		t.Fatal("Expected error for empty file, got nil")
		os.Remove(outputPath)
	}

	// Verify error message mentions empty file
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("Expected error about empty file, got: %v", err)
	}
}

// TestCompressFile_SingleCharacter tests compression with repeated character
func TestCompressFile_SingleCharacter(t *testing.T) {
	// Create file with 100 'a' characters
	testData := []byte(strings.Repeat("a", 100))
	inputPath := "test_compress_single.txt"

	err := os.WriteFile(inputPath, testData, 0644)
	if err != nil {
		t.Fatal("Failed to create test file: ", err)
	}

	defer os.Remove(inputPath)

	outputPath := "test_compress_single.hf"
	err = internal.CompressFile(inputPath, outputPath)
	if err != nil {
		t.Fatal("Compression failed: ", err)
	}

	defer os.Remove(outputPath)

	// Verify compression occurred
	stats, err := internal.GetCompressionStats(inputPath, outputPath)
	if err != nil {
		t.Fatal("Failed to get stats:", err)
	}

	if stats.CompressedSize >= stats.OriginalSize {
		t.Errorf("Single character file should compress well: compressed=%d, original=%d",
			stats.CompressedSize, stats.OriginalSize)
	}

	// Verify header
	outFile, err := os.Open(outputPath)
	if err != nil {
		t.Fatal("Failed to open output:", err)
	}
	defer outFile.Close()

	header, err := internal.ReadHeader(outFile)
	if err != nil {
		t.Fatal("Failed to read header:", err)
	}

	if header.NumChars != 1 {
		t.Errorf("Expected 1 unique char, got %d", header.NumChars)
	}

	if header.FreqTable['a'] != 100 {
		t.Errorf("Expected frequency of 'a' to be 100, got %d", header.FreqTable['a'])
	}
}

// TestCompressFile_AllDifferentChars tests worst-case: all unique characters
func TestCompressFile_AllDifferentCharacters(t *testing.T) {
	// Create string with many different strings
	inputData := []byte("abcdefghijklmnopqrstuvwxyz")
	inputPath := "test_compress_unique.txt"

	err := os.WriteFile(inputPath, inputData, 0644)
	if err != nil {
		t.Fatal("Failed to create test file: ", err)
	}
	defer os.Remove(inputPath)

	outputPath := "test_compress_unique.hf"
	err = internal.CompressFile(inputPath, outputPath)
	if err != nil {
		t.Fatal("Compression failed", err)
	}

	defer os.Remove(outputPath)

	// Verify file was created
	outputInfo, err := os.Stat(outputPath)
	if err != nil {
		t.Fatal("Output file is empty")
	}
	// With all unique chars, compression may not help much
	// (or might even increase size due to header overhead)
	// Just verify it completed successfully
	if outputInfo.Size() == 0 {
		t.Errorf("Output file is empty")
	}

	// Verify header
	outputFile, err := os.Open(outputPath)
	if err != nil {
		t.Fatal("Failed to open output: ", err)
	}
	defer outputFile.Close()

	header, err := internal.ReadHeader(outputFile)
	if err != nil {
		t.Fatal("Failed to read header: ", err)
	}

	if header.NumChars != 26 {
		t.Errorf("Expected 26 unique chars, got %d", header.NumChars)
	}
}

// TestCompressFile_NonExistentInput tests error handling for missing file
func TestCompressFile_NonExistentInput(t *testing.T) {
	inputPath := "this_file_does_not_exist.txt"
	outputPath := "output.hf"

	err := internal.CompressFile(inputPath, outputPath)
	if err == nil {
		t.Fatal("Expected error for non-existent file got nil")
		os.Remove(outputPath)
	}
}

// TestCompressFile_LargerText tests with more realistic text
func TestCompressFile_LargerText(t *testing.T) {
	// Create test data with repeated patterns (good for compression)
	testData := []byte(strings.Repeat("Hello, World! ", 50))
	inputPath := "test_compress_large.txt"

	err := os.WriteFile(inputPath, testData, 0644)
	if err != nil {
		t.Fatal("Failed to create test file: ", err)
	}
	defer os.Remove(inputPath)

	outputPath := "test_compare_larger.hf"
	err = internal.CompressFile(inputPath, outputPath)
	if err != nil {
		t.Fatal("Compression failed: ", err)
	}
	defer os.Remove(outputPath)

	// Get Comrpession sttaistics
	stats, err := internal.GetCompressionStats(inputPath, outputPath)
	if err != nil {
		t.Fatal("Failed to get stats: ", err)
	}

	// Verify stats make sense
	if stats.OriginalSize == 0 {
		t.Error("Original size should not be zero")
	}

	if stats.CompressedSize == 0 {
		t.Error("Compressed size should not be zero")
	}

	// Compression ratio should be less than 100^ for repetitive txt
	if stats.CompressionRatio > 100 {
		t.Logf("Warning: File expanded during compression (ratio: %.2f%%)", stats.CompressionRatio)
		t.Logf("This can happen with very small files due to header overhead")
	}

	t.Logf("Compression stats: Original=%d, Compressed=%d, Ratio=%.2f%%",
		stats.OriginalSize, stats.CompressedSize, stats.CompressionRatio)
}

// TestGetCompressionStats tests statistics calculation
func TestGetCompressionStats(t *testing.T) {
	// Create test files
	inputPath := "test_stats_input.txt"
	outputPath := "test_stats_output.txt"

	inputData := []byte(strings.Repeat("test", 100))
	outputData := []byte("compressed")

	err := os.WriteFile(inputPath, inputData, 0644)
	if err != nil {
		t.Fatal("Failed to create input file:", err)
	}
	defer os.Remove(inputPath)

	err = os.WriteFile(outputPath, outputData, 0644)
	if err != nil {
		t.Fatal("Failed to create output file:", err)
	}
	defer os.Remove(outputPath)

	// Get stats
	stats, err := internal.GetCompressionStats(inputPath, outputPath)
	if err != nil {
		t.Fatal("Failed to get stats:", err)
	}

	// Verify calculations
	expectedOriginal := int64(400)  // "test" * 100 = 400 bytes
	expectedCompressed := int64(10) // "compressed" = 10 bytes

	if stats.OriginalSize != expectedOriginal {
		t.Errorf("Expected original size %d, got %d", expectedOriginal, stats.OriginalSize)
	}

	if stats.CompressedSize != expectedCompressed {
		t.Errorf("Expected compressed size %d, got %d", expectedCompressed, stats.CompressedSize)
	}

	expectedRatio := (float64(10) / float64(400)) * 100 // 2.5%
	if stats.CompressionRatio != expectedRatio {
		t.Errorf("Expected ratio %.2f%%, got %.2f%%", expectedRatio, stats.CompressionRatio)
	}

	expectedSaved := int64(390)
	if stats.BytesSaved != expectedSaved {
		t.Errorf("Expected %d bytes saved, got %d", expectedSaved, stats.BytesSaved)
	}
}

// TestGetCompressionStats_NonExistentFiles tests error handling
func TestGetCompressionStats_NonExistentFiles(t *testing.T) {
	_, err := internal.GetCompressionStats("nonexistent1.txt", "nonexistent2.txt")
	if err == nil {
		t.Fatal("Expected error for non-existent files, got nil")
	}
}

// TestCompressFile_BinaryData tests with binary data (not just text)
func TestCompressFile_BinaryData(t *testing.T) {
	// Create binary data with ALL 256 possible byte values
	testData := make([]byte, 256)
	for i := 0; i < 256; i++ {
		testData[i] = byte(i)
	}

	inputPath := "test_compress_binary.bin"
	err := os.WriteFile(inputPath, testData, 0644)
	if err != nil {
		t.Fatal("Failed to create binary file:", err)
	}
	defer os.Remove(inputPath)

	outputPath := "test_compress_binary.hf"
	err = internal.CompressFile(inputPath, outputPath)

	// Should error because 256 unique bytes exceeds uint8 limit
	if err == nil {
		t.Fatal("Expected error for 256 unique bytes, got nil")
		os.Remove(outputPath)
	}

	if !strings.Contains(err.Error(), "255") {
		t.Errorf("Expected error about 255 limit, got: %v", err)
	}
}
