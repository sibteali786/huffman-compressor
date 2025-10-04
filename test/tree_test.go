package test

import (
	"huffman-compressor/internal"
	"testing"
)

func TestBuildHuffmanTree_SimpleCase(t *testing.T) {
	freqTable := internal.FrequencyTable{
		'a': 5,
		'b': 2,
		'c': 1,
	}

	// Execute
	root, err := internal.BuildHuffmanTree(freqTable)
	if err != nil {
		t.Fatalf("Failed to create Huffman Tree: %v", err)
	}

	// Verify
	if root.IsNil() {
		t.Errorf("Expected root to not be null but got %v", root)
	}

	if root.IsLeaf() {
		t.Errorf("Expected root to be an internal node")
	}

	if root.GetFreq() != 8 {
		t.Errorf("Expected root frequency to be 8 but got %d", root.GetFreq())
	}

	// 1. All original characters appear as leaves
	leaves := CollectLeaves(*root)
	if len(leaves) != len(freqTable) {
		t.Errorf("Expected %d leaves but got %d", len(freqTable), len(leaves))
	}

	charSet := make(map[byte]bool)
	for char := range freqTable {
		charSet[char] = true
	}

	for _, leaf := range leaves {
		if !leaf.IsLeaf() {
			t.Errorf("Expected leaf node but got internal node")
		}
		if _, exists := charSet[leaf.GetChar()]; !exists {
			t.Errorf("Leaf character %c not found in original frequency table", leaf.GetChar())
		}
	}

	// 2. Frequencies are correct
	for _, leaf := range leaves {
		expectedFreq, exists := freqTable[leaf.GetChar()]
		if !exists {
			t.Errorf("Leaf character %c not found in original frequency table", leaf.GetChar())
			continue
		}
		if leaf.GetFreq() != expectedFreq {
			t.Errorf("Expected frequency of character %c to be %d but got %d", leaf.GetChar(), expectedFreq, leaf.GetFreq())
		}
	}

	// 3. Tree satisfies Huffman property
	// Most frequent char should have shortest path to root
	charDepths := make(map[byte]int)
	for char := range freqTable {
		depth := GetDepth(root, char, 0)
		if depth == -1 {
			t.Errorf("Character %c not found in tree", char)
		} else {
			charDepths[char] = depth
		}
	}

	for charA, freqA := range freqTable {
		for charB, freqB := range freqTable {
			if freqA > freqB {
				if charDepths[charA] > charDepths[charB] { // Changed from >= to >
					t.Errorf("Character %c with higher frequency should have depth <= %c, got %d > %d",
						charA, charB, charDepths[charA], charDepths[charB])
				}
			}
		}
	}
}

func TestBuildHuffmanTree_EmptyTable(t *testing.T) {
	freqTable := internal.FrequencyTable{}

	// Execute
	root, err := internal.BuildHuffmanTree(freqTable)
	if err == nil {
		t.Fatalf("Expected error for empty frequency table but got none")
	}

	if !root.IsNil() {
		t.Errorf("Expected root to be nil for empty frequency table but got %v", root)
	}
}

func TestBuildHuffmanTree_TwoCharacters(t *testing.T) {
	freqTable := internal.FrequencyTable{
		'x': 10,
		'y': 5,
	}

	// Execute
	root, err := internal.BuildHuffmanTree(freqTable)
	if err != nil {
		t.Fatalf("Failed to create Huffman Tree: %v", err)
	}

	// Verify
	if root.IsNil() {
		t.Fatalf("Expected root to not be null")
	}

	if root.IsLeaf() {
		t.Errorf("Expected root to be an internal node")
	}

	if root.GetFreq() != 15 {
		t.Errorf("Expected root frequency to be 15 but got %d", root.GetFreq())
	}

	left := root.GetLeft()
	right := root.GetRight()
	if left == nil || right == nil {
		t.Fatalf("Expected both left and right children to be non-nil")
	}

	if !left.IsLeaf() || !right.IsLeaf() {
		t.Errorf("Expected both children to be leaf nodes")
	}
	if left.GetFreq() != 5 || right.GetFreq() != 10 {
		t.Errorf("Expected left freq 5 and right freq 10 but got left %d and right %d", left.GetFreq(), right.GetFreq())
	}
	if left.GetChar() != 'y' || right.GetChar() != 'x' {
		t.Errorf("Expected left char 'y' and right char 'x' but got left '%c' and right '%c'", left.GetChar(), right.GetChar())
	}

}
func CollectLeaves(node internal.HuffmanNode) []internal.HuffmanNode {
	if node.IsNil() {
		return []internal.HuffmanNode{}
	}

	if node.IsLeaf() {
		return []internal.HuffmanNode{node}
	}

	leftLeaves := CollectLeaves(*node.GetLeft())
	rightLeaves := CollectLeaves(*node.GetRight())

	return append(leftLeaves, rightLeaves...)
}
func GetDepth(node *internal.HuffmanNode, char byte, currentDepth int) int {
	if node == nil {
		return -1
	}
	if node.IsLeaf() && node.GetChar() == char {
		return currentDepth
	}

	leftDepth := GetDepth(node.GetLeft(), char, currentDepth+1)
	if leftDepth != -1 {
		return leftDepth
	}

	return GetDepth(node.GetRight(), char, currentDepth+1)
}
