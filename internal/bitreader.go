package internal

import "io"

type BitReader struct {
	reader      io.Reader // where to read bytes from
	currentByte byte      // Current byte being read
	bitPosition int       // position in current byte (0-7)
	finished    bool      // No more data to read
}

func NewBitReader(reader io.Reader) *BitReader {
	return &BitReader{
		reader:      reader,
		currentByte: 0,
		bitPosition: 8, // Start at 8 to trigger first read
		finished:    false,
	}
}

func (br *BitReader) ReadBit() (uint8, error) {
	// Reads a single bit (0 or 1)
	// If we've read all bits from current byte, get next byte
	if br.bitPosition >= 8 {
		// Read next byte from stream
		buf := make([]byte, 1)
		n, err := br.reader.Read(buf)
		if err == io.EOF || n == 0 {
			br.finished = true
			return 0, io.EOF
		}
		if err != nil {
			return 0, err
		}

		br.currentByte = buf[0]
		br.bitPosition = 0 // Reset to read from MSB
	}

	// Extract bit at current position
	// MSB is position 0, LSB is position 7
	bitValue := (br.currentByte >> (7 - br.bitPosition)) & 1
	br.bitPosition++

	return bitValue, nil
}

func (br *BitReader) IsFinished() bool {
	return br.finished
}
