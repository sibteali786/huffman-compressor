package internal

import (
	"fmt"
	"unicode"
)

type HuffmanCode struct {
	bits   uint64
	length int
}
type CodeTable = map[byte]HuffmanCode

func (hfc *HuffmanCode) GetLength() int {
	return hfc.length
}

func (hfc *HuffmanCode) GetBits() int {
	return int(hfc.bits)
}

func GenerateCodes(root *HuffmanNode) CodeTable {
	codeTable := make(CodeTable)

	//case: root == nil
	if root == nil {
		return nil
	}

	// case: single character
	leftIsDummy := root.left != nil && root.left.IsDummy()
	rightIsDummy := root.right != nil && root.right.IsDummy()

	if leftIsDummy && !rightIsDummy && root.right.IsLeaf() {
		codeTable[root.right.char] = HuffmanCode{bits: 1, length: 1}
		return codeTable
	} else if rightIsDummy && !leftIsDummy && root.left.IsLeaf() {
		codeTable[root.left.char] = HuffmanCode{bits: 0, length: 1}
		return codeTable
	}
	// normal case: traverse and build codes
	var initialCode HuffmanCode
	buildCodesRecursive(root, initialCode, codeTable)
	return codeTable
}

func buildCodesRecursive(node *HuffmanNode, currentCode HuffmanCode, codeTable CodeTable) {
	if node == nil {
		return
	}

	if node.IsLeaf() {
		// Found a char store its code
		codeTable[node.char] = currentCode
		return
	}
	// Recurse left ( append 0 )
	leftCode := appendBit(currentCode, 0)
	buildCodesRecursive(node.left, leftCode, codeTable)

	// Recurse right ( append 1 )
	rightCode := appendBit(currentCode, 1)
	buildCodesRecursive(node.right, rightCode, codeTable)
}

func PrintCodeTable(codeTable CodeTable) {
	for char, code := range codeTable {
		switch char {
		case ' ':
			fmt.Printf("Space: %v\n", code)
		case '\n':
			fmt.Printf("Newline: %v\n", code)
		case '\t':
			fmt.Printf("Tab: %v\n", code)
		default:
			if unicode.IsPrint(rune(char)) {
				fmt.Printf("' %v ': %v\n", string(char), code)
			} else {
				fmt.Printf("(0x%02x): %v\n", char, CodeToString(code))
			}

		}
	}
}

func appendBit(code HuffmanCode, bit uint8) HuffmanCode {
	newCode := code
	if bit != 0 {
		// if bit is 1 we left shift 0...001 by code.length so it becomes 0...(code.length)00.. and then bitwiseOR with existing bits
		// essentially its like adding 1 to third position from left of 000011
		newCode.bits |= (1 << code.length)

	}
	// if bit is 0 we don't need to do anything as the bit is already 0
	newCode.length++
	return newCode
}

func CodeToString(code HuffmanCode) string {
	if code.length == 0 {
		return ""
	}

	result := make([]byte, code.length)
	for i := uint8(0); i < uint8(code.length); i++ {
		// Read bits from left to right (MSB to LSB for display)a
		// We read from position (length - 1 - i) to get MSB first
		bitPos := uint8(code.length) - 1 - i
		if (code.bits>>bitPos)&1 == 1 {
			result[i] = '1'
		} else {
			result[i] = '0'
		}
	}

	return string(result)
}

func VerifyPrefixFree(codeTable CodeTable) bool {
	// for each pair of codes, verify neither is prefix of other
	codes := make([]HuffmanCode, 0, len(codeTable))
	for _, code := range codeTable {
		codes = append(codes, code)
	}

	for i := 0; i < len(codes); i++ {
		for j := i + 1; j < len(codes); j++ {
			if isPrefix(codes[i], codes[j]) || isPrefix(codes[j], codes[i]) {
				return false
			}
		}
	}
	return true
}

func isPrefix(code1, code2 HuffmanCode) bool {
	if code1.length >= code2.length {
		return false
	}

	// Both codes store bits LSB-first (position 0 = first bit appended)
	// To check if code1 is a prefix, we compare the first code1.length bits

	// Extract the first code1.length bits from code2
	// Since both are stored the same way, we just mask the lower bits
	mask := uint64((1 << code1.length) - 1)
	code2Prefix := code2.bits & mask

	return code2Prefix == code1.bits
}
