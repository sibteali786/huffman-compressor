package internal

import (
	"fmt"
	"sort"
)

type HuffmanNode struct {
	char      byte // the character (only for leaf nodes)
	frequency int
	left      *HuffmanNode
	right     *HuffmanNode
	isLeaf    bool
}

func (hf *HuffmanNode) IsLeaf() bool {
	return hf.isLeaf
}

func (hf *HuffmanNode) GetFreq() int {
	return hf.frequency
}

func (hf *HuffmanNode) IsNil() bool {
	return hf == nil
}

func (hf *HuffmanNode) GetChar() byte {
	return hf.char
}

func (hf *HuffmanNode) GetLeft() *HuffmanNode {
	return hf.left
}

func (hf *HuffmanNode) GetRight() *HuffmanNode {
	return hf.right
}

func (hf *HuffmanNode) IsDummy() bool {
	return hf != nil && hf.isLeaf && hf.char == 0 && hf.frequency == 0
}

func NewLeafNode(char byte, frequency int) *HuffmanNode {
	return &HuffmanNode{
		char:      char,
		frequency: frequency,
		isLeaf:    true,
		left:      nil, right: nil,
	}
}

func NewInternalNode(left *HuffmanNode, right *HuffmanNode) *HuffmanNode {
	return &HuffmanNode{
		frequency: left.frequency + right.frequency,
		isLeaf:    false,
		left:      left,
		right:     right,
	}
}

func BuildHuffmanTree(freqTable FrequencyTable) (*HuffmanNode, error) {
	// Edge case: empty table
	if len(freqTable) == 0 {
		return nil, fmt.Errorf("frequency table has no entries to process")
	}

	// ✅ Edge case: only one character - COMPLETE IMPLEMENTATION
	if len(freqTable) == 1 {
		var singleChar byte
		var singleFreq int
		for char, freq := range freqTable {
			singleChar = char
			singleFreq = freq
			break
		}

		// Create a leaf node for the character
		leafNode := NewLeafNode(singleChar, singleFreq)

		// Create a dummy internal node
		dummyNode := NewLeafNode(0, 0)

		// Create a root with real leaf on left and dummy on right
		root := NewInternalNode(leafNode, dummyNode)
		return root, nil
	}

	// Calculate total unique characters for Priority Queue
	capacity := len(freqTable)

	// Create priority queue
	pq := NewPriorityQueue[*HuffmanNode](capacity * 2)

	// ✅ Sort characters to ensure deterministic tree building
	chars := make([]byte, 0, len(freqTable))
	for char := range freqTable {
		chars = append(chars, char)
	}

	// Sort characters (ensures same order every time)
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})

	// Step 1: Create leaf nodes and enqueue (in sorted order)
	for _, char := range chars {
		freq := freqTable[char]
		leafNode := NewLeafNode(char, freq)
		pq.Enqueue(leafNode, freq)
	}

	// Step 2: Build tree bottom up
	for pq.Size() > 1 {
		left, errLeft := pq.Dequeue()
		right, errRight := pq.Dequeue()
		if errLeft != nil || errRight != nil {
			return nil, fmt.Errorf("error getting nodes from PQ: %v, %v", errLeft, errRight)
		}
		parentNode := NewInternalNode(left, right)
		pq.Enqueue(parentNode, parentNode.frequency)
	}

	// Step 3: root is last remaining node
	root, err := pq.Dequeue()
	if err != nil {
		return nil, fmt.Errorf("error getting root node: %v", err)
	}
	return root, nil
}

func PrintTree(node *HuffmanNode, prefix string, isLeft bool) {
	if node == nil {
		return
	}
	var newPrefix string
	if isLeft {
		fmt.Println(prefix)
		fmt.Println("├──")
		newPrefix = prefix + "│   "
	} else {
		fmt.Println("└──")
		newPrefix = prefix + "    "
	}

	// print node information
	if node.isLeaf {
		// For leaf: show character and frequency
		fmt.Printf("Leaf: ' %v ' (freq:  %v )", node.char, node.frequency)
	} else {
		fmt.Printf("Internal (freq: ' %v ')", node.frequency)
	}

	fmt.Println()

	// recursively print children
	if node.left != nil {
		PrintTree(node.left, newPrefix, true)
	}

	if node.right != nil {
		PrintTree(node.right, newPrefix, false)
	}
}
