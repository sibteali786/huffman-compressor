package internal

import (
	"io"
)

type BitBuffer struct {
	currentByte byte      // Current byte being built
	bitPosition int       // Position in current byte (0-7)
	writer      io.Writer // Where to write complete bytes
	totalBits   int       // For statistics
}

func NewBitBuffer(writer io.Writer) *BitBuffer {
	return &BitBuffer{
		currentByte: 0,
		bitPosition: 0,
		writer:      writer,
		totalBits:   0,
	}
}

// WriteBit writes a single bit and flushes when byte is full
func (bb *BitBuffer) WriteBit(bit uint64) {
	// Set bit in current byte
	if bit == 1 {
		mask := byte(1 << (7 - bb.bitPosition))
		bb.currentByte |= mask
	}

	bb.bitPosition++
	bb.totalBits++

	// if byte is complete, write it to file immediately
	if bb.bitPosition == 8 {
		bb.Flush()
	}

}

func (bb *BitBuffer) WriteBits(bits uint64, length int) {
	// Write bits from LSB (position 0) to MSB (position length-1)
	// This matches how appendBit() stores them
	for i := 0; i < length; i++ {
		bit := (bits >> i) & 1
		bb.WriteBit(bit)
	}
}

func (bb *BitBuffer) Flush() {
	// Write byte to file
	bb.writer.Write([]byte{bb.currentByte})

	// Reset for next byte
	bb.currentByte = 0
	bb.bitPosition = 0
}

func (bb *BitBuffer) Close() (int, error) {
	paddingBits := bb.GetPaddingBits()
	// if there are incomplete bits, pad and flush
	if bb.bitPosition > 0 {
		// currentByte already has the bits, just write it
		bb.writer.Write([]byte{bb.currentByte})
	}

	return paddingBits, nil
}

func (bb *BitBuffer) GetTotalBits() int {
	return bb.totalBits
}

func (bb *BitBuffer) GetPaddingBits() int {
	if bb.bitPosition == 0 {
		return 0
	}
	return 8 - bb.bitPosition
}
