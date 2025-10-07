package test

import (
	"huffman-compressor/internal"
	"testing"
)

func TestGenerateCodes_SimpleCase(t *testing.T) {
	freqTable := internal.FrequencyTable{'a': 5,
		'b': 2,
		'c': 1,
	}

	root, err := internal.BuildHuffmanTree(freqTable)
	if err != nil {
		t.Fatalf("There is some problem: %v", err)
	}

	if root == nil {
		t.Fatal(err)
	}

	codeTable := internal.GenerateCodes(root)
	// Verify
	if len(codeTable) != 3 {
		t.Fatalf("Expected length of codeTable to be %v but got %v", len(freqTable), len(codeTable))
	}

	for char := range freqTable {
		if _, exists := codeTable[char]; !exists {
			t.Fatalf("Expected codeTable to contain %v, but didn't found it", char)
		}
	}

	// verify huffman property
	for charA, hfCodeA := range codeTable {
		for charB, hfCodeB := range codeTable {
			if charA == charB {
				continue
			}

			freqA := freqTable[charA]
			freqB := freqTable[charB]

			lenA := hfCodeA.GetLength()
			lenB := hfCodeB.GetLength()

			if freqA > freqB && lenA > lenB {
				t.Errorf("'%c' (freq=%d, len=%d) should have ≤ code length than '%c' (freq=%d, len=%d)",
					charA, freqA, lenA, charB, freqB, lenB)
			}
		}
	}
	// Verify Prefix free
	if internal.VerifyPrefixFree(codeTable) != true {
		t.Fatal("The pairs created are not prefix free")
	}

	// verify no code is empty
	for char, code := range codeTable {
		if code.GetLength() == 0 {
			t.Fatalf("Expected length of %v's code to be non zero but got %v", char, code.GetLength())
		}
	}

}

func TestGenerateCodes_SingleChar(t *testing.T) {
	freqTable := internal.FrequencyTable{
		'x': 100,
	}

	root, err := internal.BuildHuffmanTree(freqTable)
	if err != nil {
		t.Fatal(err)
	}

	if root == nil {
		t.Fatal("Root is nill")
	}

	codeTable := internal.GenerateCodes(root)

	// assert lengths of freq table and code table are same
	if len(codeTable) != 1 {
		t.Fatalf("Expected length of codeTable to be %v got %v", len(freqTable), len(codeTable))
	}
	// assert: codeTable contains x
	for char := range freqTable {
		if _, exists := codeTable[char]; !exists {
			t.Fatalf("Expected codeTable to contain %v, but didn't found it", char)
		}
	}
	// assert len of codeTable['x'] is > 0
	for char, code := range codeTable {
		if code.GetLength() == 0 {
			t.Fatalf("Expected length of %v's code to be non zero but got %v", char, code.GetLength())
		}
	}

}
func TestGenerateCodes_TwoChars(t *testing.T) {
	freqTable := internal.FrequencyTable{
		'a': 10,
		'b': 5,
	}

	root, err := internal.BuildHuffmanTree(freqTable)
	if err != nil {
		t.Fatal(err)
	}

	if root == nil {
		t.Fatal("Root is nill")
	}

	codeTable := internal.GenerateCodes(root)

	// assert lengths of freq table and code table are same
	if len(codeTable) != 2 {
		t.Fatalf("Expected length of codeTable to be %v got %v", len(freqTable), len(codeTable))
	}
	// assert: codeTable contains x
	for charA, hfCodeA := range codeTable {
		for charB, hfCodeB := range codeTable {
			if charA == charB {
				continue
			}

			freqA := freqTable[charA]
			freqB := freqTable[charB]

			lenA := hfCodeA.GetLength()
			lenB := hfCodeB.GetLength()

			if freqA > freqB && lenA > lenB {
				t.Errorf("'%c' (freq=%d, len=%d) should have ≤ code length than '%c' (freq=%d, len=%d)",
					charA, freqA, lenA, charB, freqB, lenB)
			}
		}
	}
	charA := codeTable['a']
	charB := codeTable['b']

	if charA.GetBits() == charB.GetBits() {
		t.Fatal("The code length of 'a' and 'b' is not equal")
	}
}

func TestGenerateCodes_EmptyTree(t *testing.T) {
	codeTable := internal.GenerateCodes(nil)

	if codeTable != nil {
		t.Fatal("CodeTable is not nill")
	}

	if len(codeTable) != 0 {
		t.Fatal("The codeTable is not zero")
	}
}
