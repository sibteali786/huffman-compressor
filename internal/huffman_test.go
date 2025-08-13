package internal

import (
	"os"
	"testing"
)

func TestAnalyzeFrequencies(t *testing.T) {
	testContent := "aaabbc"
	testFile := "test_input.txt"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	defer os.Remove(testFile)

	freqTable, err := AnalyzeFrequencies(testFile)
	if err != nil {
		t.Fatalf("AnalyzeFrequencies failed: %v", err)
	}

	expected := map[byte]int{
		'a': 3,
		'b': 2,
		'c': 1,
	}

	if len(freqTable) != len(expected) {
		t.Errorf("Expected %d unique characters, got %d", len(expected), len(freqTable))
	}

	for char, expectedCount := range expected {
		if actualCount, exists := freqTable[char]; !exists {
			t.Errorf("Character '%c' not found in frequency table", char)
		} else if actualCount != expectedCount {
			t.Errorf("Character '%c': expected %d, got %d", char, expectedCount, actualCount)
		}
	}
}
