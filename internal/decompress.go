package internal

import (
	"fmt"
	"io"
	"os"
)

func Decompress(inputPath, outputPath string) error {
	// ==================== PHASE 1: Open Compressed File ====================
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open compressed file: %s", err)
	}
	defer inputFile.Close()

	// ==================== PHASE 2: Read and parse Header ====================
	header, err := ReadHeader(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read header: %s", err)
	}
	// Extract info from header
	originalSize := header.OriginalSize
	freqTable := header.FreqTable
	// paddingBits := header.PaddingBits

	// Validate header
	if originalSize == 0 {
		return fmt.Errorf("invalid header: original size is 0")
	}

	if len(freqTable) == 0 {
		return fmt.Errorf("invalid header: freq table is empty")
	}

	// ==================== PHASE 3: Rebuild Huffman Tree ====================
	// Rebuild the tree from frequencies (same as compression)
	root, err := BuildHuffmanTree(freqTable)
	if err != nil {
		return fmt.Errorf("failed to create output file: %s", err)
	}

	// ==================== PHASE 4: Create Output File ====================

	// Create output file for decompressed data
	outputFile, err := os.Create(outputPath)
	if err != nil {

		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer outputFile.Close()

	// ==================== PHASE 5: Decode Bit Stream ====================
	// Create bit reader for the compressed data
	// (inputFile is already positioned after header)
	bitReader := NewBitReader(inputFile)

	// Buffer for writing decoded bytes for efficiency
	writeBuffer := make([]byte, 0, 1024)

	// track how many bytes we decoded
	bytesDecoded := uint64(0)

	// Start at root of tree
	currentNode := root

	// Read bits until we decoded all bytes
	for bytesDecoded < originalSize {
		// Read next bit
		bit, err := bitReader.ReadBit()

		// Handle EOF or errors
		if err == io.EOF {
			// Reached end of file
			// This is expected when we've decoded all bytes
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read bit: %s", err)
		}

		// Traverse tree based on bit value
		if bit == 0 {
			currentNode = currentNode.left
		} else {
			currentNode = currentNode.right
		}

		// Safety check: if not at leaf, must have children
		if !currentNode.IsLeaf() && currentNode.left == nil {
			return fmt.Errorf("corrupted data: invalid tree traversal")
		}
		// Check if we hit a leaf node
		if currentNode.IsLeaf() {
			// Found a character
			char := currentNode.GetChar()

			// Add to write buffer
			writeBuffer = append(writeBuffer, char)
			bytesDecoded++

			// Flush buffer when it gets large (for efficiency)
			if len(writeBuffer) >= 1024 {
				_, err := outputFile.Write(writeBuffer)
				if err != nil {
					return fmt.Errorf("failed to write output: %s", err)
				}
				writeBuffer = writeBuffer[:0] // Reset Buffer
			}
			// Reset to root for next character
			currentNode = root
		}
	}

	// ==================== PHASE 6: Flush Remaining Data ====================
	/// Write any remaining bytes in buffer
	if len(writeBuffer) > 0 {
		_, err := outputFile.Write(writeBuffer)
		if err != nil {
			return fmt.Errorf("failed to write final output: %s", err)
		}
	}

	// ==================== PHASE 7: Validate Result ====================

	// Verify we decoded the correct number of bytes
	if bytesDecoded != originalSize {
		return fmt.Errorf("decoded %d bytes, expected %d", bytesDecoded, originalSize)
	}
	return nil
}

func VerifyDecompression(originalPath, decompressedPath string) error {
	// Read both files
	originalData, err := os.ReadFile(originalPath)
	if err != nil {
		return err
	}

	decompressedData, err := os.ReadFile(decompressedPath)

	if err != nil {
		return err
	}

	// compare sizes
	if len(originalData) != len(decompressedData) {
		return fmt.Errorf("size mismatch: original=%d, decompress=%d", len(originalData), len(decompressedData))
	}

	//Compare byte to byte
	for i := 0; i < len(originalData); i++ {
		if originalData[i] != decompressedData[i] {
			return fmt.Errorf("mismatch at byte %d: original=%d, decompressed=%d", i, originalData[i], decompressedData[i])
		}
	}

	return nil
}
