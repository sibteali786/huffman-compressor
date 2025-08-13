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

		switch key {
		case ' ':
			fmt.Printf("Space: %v occurrences\n", value)
		case '\n':
			fmt.Printf("Newline: %v occurrences\n", value)
		case '\t':
			fmt.Printf("Tab: %v occurrences\n", value)
		default:
			// Only print alphanumeric characters
			if unicode.IsLetter(rune(key)) || unicode.IsDigit(rune(key)) {
				fmt.Printf("Character '%v' (%v): %v occurrences\n", string(key), key, value)
			} else if unicode.IsPrint(rune(key)) {
				fmt.Printf("Character '%c' (%d): %d occurrences\n", key, key, value)
			} else {
				fmt.Printf("Non-printable character (0x%02x): %d occurrences\n", key, value)
			}
		}

	}
}
