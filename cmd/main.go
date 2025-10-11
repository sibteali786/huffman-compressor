package main

import (
	"flag"
	"fmt"
	"huffman-compressor/internal"
	"os"
)

func main() {
	var (
		inputFile  = flag.String("input", "", "Input file to compress/decompress")
		outputFile = flag.String("output", "", "Output file")
		compress   = flag.Bool("compress", false, "Compress the input file")
		decompress = flag.Bool("decompress", false, "Decompress the input file")
	)

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: Input file is required")
		flag.Usage()
		os.Exit(1)
	}
	// make sure only one of these is set and not both
	if *compress && *decompress {
		fmt.Println("Error: Provide either compress or decompress option, not both")
		flag.Usage()
		os.Exit(1)
	}
	// if not output file name provided then
	if *outputFile == "" {
		*outputFile = "output.txt"
	}

	// validate input file exists
	_, err := os.Stat(*inputFile)

	if os.IsNotExist(err) {
		fmt.Printf("%v file does not exists\n", inputFile)
		os.Exit(1)
	}

	table, err := internal.AnalyzeFrequencies(*inputFile)

	if err != nil {
		fmt.Printf("Error analyzing frequencies: %v\n", err)
		os.Exit(1)
	}

	internal.PrintFrequencies(table)

	if *compress {
		fmt.Printf("Compressing %s to %s...\n", *inputFile, *outputFile)

		// Perform compression
		err := internal.CompressFile(*inputFile, *outputFile)
		if err != nil {

			fmt.Fprintf(os.Stderr, "Compression failed: %v\n", err)
			os.Exit(1)
		}

		// Show statistics
		stats, err := internal.GetCompressionStats(*inputFile, *outputFile)
		if err == nil {

			internal.PrintCompressionStats(stats)
		}

		fmt.Printf("✓ Successfully compressed to %s\n", *outputFile)
	} else if *decompress {
		fmt.Printf("Decompressing %s to %s\n", *inputFile, *outputFile)
		// TODO: Will implement in Step 6-7
		fmt.Println("Decompression not yet implemented")
	}
}
