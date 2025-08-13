package internal

import (
	"fmt"
	"io"
	"os"
	"unicode"
)

type FrequencyTable map[byte]int

func AnalyzeFrequencies(filename string) (FrequencyTable, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	table := make(FrequencyTable)

	for {
		count, err := file.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			table[buffer[i]]++
		}
	}
	return table, nil
}

func PrintFrequencies(freqTable FrequencyTable) {
	// Your task: Implement this function
	// Print the frequency table in a readable format
	// Example output: "Character 'a' (97): 1234 occurrences"
	// Handle special characters (spaces, newlines) nicely
	for key, value := range freqTable {
		if unicode.IsLetter(rune(key)) || unicode.IsDigit(rune(key)) {
			fmt.Printf("Character '%v' (%v): %v occurrences\n", string(key), key, value)
		}
	}
}
