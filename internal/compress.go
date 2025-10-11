package internal

import (
	"fmt"
	"io"
	"os"
)

func CompressFile(inputPath, outputPath string) error {
	// Use existing streaming frequency analysis
	freqTable, err := AnalyzeFrequencies(inputPath)
	if err != nil {
		return fmt.Errorf("failed to analyze frequencies: %w", err)
	}

	// Edge case: empty file (frequency table will be empty)
	if len(freqTable) == 0 {
		return fmt.Errorf("cannot compress empty file")
	}

	// Get original file size for header
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	originalSize := uint64(fileInfo.Size())

	// ==================== PHASE 2: Build Tree and Codes ====================

	// Build Huffman Tree
	root, err := BuildHuffmanTree(freqTable)
	if err != nil {
		return fmt.Errorf("failed to build huffman tree: %s", err)
	}

	if len(freqTable) > 255 {
		return fmt.Errorf("file contains %d unique byte values (max 255 supported)", len(freqTable))
	}

	// Generate code table from tree
	codeTable := GenerateCodes(root)

	// Optional: verify codes are prefix-free (debug mode)
	if !VerifyPrefixFree(codeTable) {
		return fmt.Errorf("generated codes are not prefix-free")
	}

	// ==================== PHASE 3: Create Output File ====================

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create the output file")
	}

	defer outputFile.Close()

	// ==================== PHASE 4: Write Header (Placeholder) ====================
	// Write header with padding = 0 (we'll update this later)
	err = WriteHeader(outputFile, freqTable, originalSize, 0)
	if err != nil {
		return fmt.Errorf("failed to write header: %s", err)
	}

	// ==================== PHASE 5: Encode and Write Data ====================
	// Open input file for reading
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %s", err)
	}
	defer inputFile.Close()
	// Create bit buffer that writes to the file
	bitBuffer := NewBitBuffer(outputFile)

	// Read and encode file in chunks (streaming approach)
	buffer := make([]byte, 1024)
	// Encode each byte in the input data
	for {
		count, err := inputFile.Read(buffer)

		// Encode each byte in this chunk
		for i := 0; i < count; i++ {
			char := buffer[i]
			code := codeTable[char]
			bitBuffer.WriteBits(code.bits, code.length)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
	}

	// Get Padding bits before closing
	paddingBits := bitBuffer.GetPaddingBits()

	// Close Buffer (flushes remaining bits with padding)
	_, err = bitBuffer.Close()
	if err != nil {
		return fmt.Errorf("failed to close bit buffer: %s", err)
	}

	// ==================== PHASE 6: Update Padding in Header ====================

	// Seek to padding byte position in header
	// Header structure: [HF:2][OrigSize:8][NumChars:1][PaddingBits:1][FreqEntries...]
	// Padding byte is at offset 11 (2 + 8 + 1)

	paddingByteOffset := int64(11)

	_, err = outputFile.Seek(paddingByteOffset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek to padding byte: %s", err)
	}

	// Write the actual padding value
	_, err = outputFile.Write([]byte{uint8(paddingBits)})
	if err != nil {
		return fmt.Errorf("failed to update padding byte: %s", err)
	}

	return nil
}

type CompressionStats struct {
	OriginalSize     int64
	CompressedSize   int64
	CompressionRatio float64
	BytesSaved       int64
}

func GetCompressionStats(inputPath string, outputPath string) (CompressionStats, error) {

	stats := CompressionStats{}

	// Get input file size
	inputInfo, err := os.Stat(inputPath)
	if err != nil {

		return stats, err
	}

	// Get output file size
	outputInfo, err := os.Stat(outputPath)
	if err != nil {

		return stats, err
	}

	stats.OriginalSize = inputInfo.Size()
	stats.CompressedSize = outputInfo.Size()

	// Calculate compression ratio (as percentage)
	stats.CompressionRatio = (float64(stats.CompressedSize) / float64(stats.OriginalSize)) * 100

	// Calculate bytes saved
	stats.BytesSaved = stats.OriginalSize - stats.CompressedSize

	return stats, nil

}

// PrintCompressionStats displays compression results
func PrintCompressionStats(stats CompressionStats) {

	fmt.Printf("\n=== Compression Statistics ===\n")
	fmt.Printf("Original size:    %d bytes\n", stats.OriginalSize)
	fmt.Printf("Compressed size:  %d bytes\n", stats.CompressedSize)
	fmt.Printf("Compression ratio: %.2f%%\n", stats.CompressionRatio)

	if stats.BytesSaved > 0 {

		fmt.Printf("Space saved:      %d bytes\n", stats.BytesSaved)
	} else {

		fmt.Printf("Space increased:  %d bytes (file too small to compress)\n", -stats.BytesSaved)
	}

	fmt.Printf("==============================\n")
}
