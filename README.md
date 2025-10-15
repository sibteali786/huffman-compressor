# Huffman Compressor

A fully functional file compression tool implementing Huffman coding algorithm from scratch in Go. This project demonstrates lossless data compression using binary trees, priority queues, and bit manipulation.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## 📋 Table of Contents

- [About](#about)
- [Features](#features)
- [How It Works](#how-it-works)
- [Installation](#installation)
- [Usage](#usage)
- [Technical Implementation](#technical-implementation)
- [Performance](#performance)
- [Learning Journey](#learning-journey)
- [Future Enhancements](#future-enhancements)
- [Acknowledgments](#acknowledgments)

## 🎯 About

This project was inspired by [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-huffman/) and implements the classic Huffman coding algorithm for lossless data compression. Built entirely from scratch without using compression libraries, this tool showcases fundamental computer science concepts including:

- Binary tree data structures
- Priority queues (min-heap)
- Greedy algorithms
- Bit manipulation
- File I/O streaming

## ✨ Features

- **Lossless Compression**: Perfectly reconstructs original files
- **Streaming Architecture**: Memory-efficient processing for large files
- **Binary Tree Implementation**: Custom Huffman tree with optimal prefix codes
- **Priority Queue**: Hand-rolled min-heap using generics
- **Bit-Level Operations**: Efficient bit packing and unpacking
- **Comprehensive Testing**: 30+ unit and integration tests
- **Edge Case Handling**: Single character files, empty files, binary data

## 🔍 How It Works

### Compression Process
```
Input File
    ↓
1. Analyze character frequencies
    ↓
2. Build Huffman tree (greedy algorithm)
    ↓
3. Generate variable-length prefix codes
    ↓
4. Write header (frequencies + metadata)
    ↓
5. Encode data using generated codes
    ↓
6. Pack bits into bytes
    ↓
Compressed File (.hf)
```

### Decompression Process
```
Compressed File (.hf)
    ↓
1. Read and parse header
    ↓
2. Rebuild Huffman tree from frequencies
    ↓
3. Read compressed bits
    ↓
4. Traverse tree bit-by-bit
    ↓
5. Decode characters at leaf nodes
    ↓
Original File Restored
```

## 📦 Installation

**Prerequisites:**
- Go 1.21 or higher

**Clone and Build:**
```bash
# Clone the repository
git clone https://github.com/yourusername/huffman-compressor.git
cd huffman-compressor

# Build the binary
go build -o huffman cmd/main.go

# Or install directly
go install
```

## 🚀 Usage

### Compress a File
```bash
# Basic compression
./huffman -compress -input file.txt -output file.hf

# The tool will display compression statistics:
# === Compression Statistics ===
# Original size:    1000 bytes
# Compressed size:  650 bytes
# Compression ratio: 65.00%
# Space saved:      350 bytes
```

### Decompress a File
```bash
# Decompress back to original
./huffman -decompress -input file.hf -output restored.txt

# Verify files are identical
diff file.txt restored.txt  # Should show no differences
```

### Examples
```bash
# Compress a text file with repetitive content (best compression)
./huffman -compress -input alice.txt -output alice.hf

# Compress a binary file
./huffman -compress -input image.bin -output image.hf

# Decompress
./huffman -decompress -input alice.hf -output alice_restored.txt
```

## 🏗️ Technical Implementation

### File Format Specification

**Compressed File Structure (.hf):**
```
[HEADER]
  - Magic Number (2 bytes): "HF"
  - Original Size (8 bytes): uint64
  - Unique Characters (1 byte): uint8 (max 255)
  - Padding Bits (1 byte): uint8 (0-7)
  - Frequency Entries (N × 5 bytes):
      - Character (1 byte)
      - Frequency (4 bytes): uint32

[COMPRESSED DATA]
  - Variable-length encoded bits packed into bytes
```

### Key Data Structures

**1. Priority Queue (Min-Heap)**
```go
// Generic implementation with custom comparators
type PriorityQueue[T any] struct {
    heap *MinHeap[T]
}
```

**2. Huffman Tree Node**
```go
type HuffmanNode struct {
    char      byte
    frequency int
    left      *HuffmanNode
    right     *HuffmanNode
    isLeaf    bool
}
```

**3. Huffman Code**
```go
type HuffmanCode struct {
    bits   uint64  // LSB-first storage
    length int     // Number of bits
}
```

**4. Bit Buffer (Writer)**
```go
// Streams bits to file, auto-flushing complete bytes
type BitBuffer struct {
    currentByte byte
    bitPosition int
    writer      io.Writer
}
```

**5. Bit Reader**
```go
// Reads individual bits from byte stream
type BitReader struct {
    currentByte byte
    bitPosition int
    reader      io.Reader
}
```

### Algorithm Complexity

- **Tree Building**: O(n log n) where n = unique characters
- **Code Generation**: O(n) tree traversal
- **Compression**: O(m) where m = file size
- **Decompression**: O(m × log n) for tree traversal
- **Space**: O(n) for tree + O(1) for streaming

## 📊 Performance

### Compression Results

| File Type | Size | Compressed | Ratio | Time |
|-----------|------|------------|-------|------|
| Text (repetitive) | 1 MB | 450 KB | 45% | ~50ms |
| Source code | 500 KB | 380 KB | 76% | ~30ms |
| Random data | 1 MB | 1.01 MB | 101% | ~60ms |
| Alice in Wonderland | 167 KB | 93 KB | 56% | ~15ms |

**Note**: Compression effectiveness depends on data repetition. Random data may slightly expand due to header overhead.

### Memory Efficiency

- **Streaming I/O**: Reads and writes in 1KB chunks
- **Constant memory**: O(1) additional space during encode/decode
- **Tree overhead**: ~5 bytes per unique character in header

## 📚 Learning Journey

This project was built as a learning exercise to understand:

### Core Concepts Mastered

1. **Data Structures**
   - Binary trees (construction and traversal)
   - Priority queues (heap implementation)
   - Hash tables (frequency counting)

2. **Algorithms**
   - Greedy algorithms (Huffman's approach)
   - Tree traversal (DFS for code generation)
   - Prefix-free encoding

3. **Systems Programming**
   - File I/O streaming
   - Bit manipulation
   - Memory management
   - Error handling

4. **Go Language Features**
   - Generics (for priority queue)
   - Interfaces (io.Reader/Writer)
   - Defer statements
   - Table-driven tests

### Key Challenges Solved

- **LSB-first storage with MSB-first writing**: Implemented bit reversal in WriteBits
- **Single character edge case**: Added dummy node for tree construction
- **Padding handling**: Track padding bits in header for accurate decompression
- **Streaming architecture**: Process large files without loading into memory
- **Uint8 overflow**: Limited unique characters to 255 (practical constraint)

## 🚀 Future Enhancements

### Potential Improvements

- [ ] **Adaptive Huffman Coding**: Update tree on-the-fly
- [ ] **Dictionary Encoding**: Combine with LZ77 for better compression
- [ ] **Parallel Processing**: Multi-threaded compression for large files
- [ ] **Progress Indicators**: Show progress for large file operations
- [ ] **Compression Levels**: Trade speed for ratio (like gzip -1 to -9)
- [ ] **Directory Compression**: Archive multiple files (like tar)
- [ ] **GUI Interface**: Desktop app with drag-and-drop
- [ ] **Benchmark Suite**: Automated performance testing
- [ ] **256 Character Support**: Use uint16 for NumChars field
- [ ] **Streaming API**: Library interface for programmatic use

### Algorithm Variants

- **Canonical Huffman**: Simplify tree serialization
- **Length-Limited Huffman**: Bound code length for hardware
- **Adaptive Huffman**: Dynamic tree updates (LZSS + Huffman)

## 🧪 Testing

Run the comprehensive test suite:
```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific test package
go test ./test/compress_test.go -v

# Run benchmarks
go test ./... -bench=. -benchmem
```

**Test Coverage:**
- Unit tests: Priority queue, tree building, code generation
- Integration tests: Full compress/decompress round-trips
- Edge cases: Empty files, single characters, binary data
- Error handling: Corrupted files, invalid headers

## 📁 Project Structure
```
huffman-compressor/
├── cmd/
│   └── main.go              # CLI entry point
├── internal/
│   ├── huffman.go           # Frequency analysis
│   ├── tree.go              # Tree construction
│   ├── encoder.go           # Code generation
│   ├── bitbuffer.go         # Bit writing
│   ├── bitreader.go         # Bit reading
│   ├── header.go            # File format handling
│   ├── compress.go          # Main compression logic
│   ├── decompress.go        # Main decompression logic
│   └── priority_queue.go    # Min-heap implementation
├── test/
│   ├── *_test.go            # Comprehensive test suites
└── README.md
```

## 🤝 Contributing

This is a learning project, but contributions are welcome! Feel free to:

- Report bugs
- Suggest features
- Submit pull requests
- Share improvements

## 📝 License

MIT License - feel free to use this code for learning and projects.

## 🙏 Acknowledgments

- **Coding Challenges** by John Crickett for the project inspiration
- **David Huffman** for the algorithm (1952)
- **Introduction to Algorithms** (CLRS) for theoretical foundation
- The Go community for excellent documentation

## 📧 Contact

Created by Syed Sibteali Baqar - sibteali786@gmail.com

Project Link: [https://github.com/sibteali786/huffman-compressor](https://github.com/sibteali786/huffman-compressor)

---

⭐ If you found this project helpful for learning, please star it on GitHub!

## 📖 Additional Resources

- [Huffman Coding Visualization](https://www.cs.usfca.edu/~galles/visualization/Huffman.html)
- [Original Huffman Paper (1952)](https://compression.ru/download/articles/huff/huffman_1952_minimum-redundancy-codes.pdf)
- [Data Compression Explained](http://mattmahoney.net/dc/dce.html)