package main

import (
	"flag"
	"fmt"
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

	if *compress {
		fmt.Printf("Compressing %s to %s\n", *inputFile, *outputFile)
		// TODO: Call compression function
	} else if *decompress {
		fmt.Printf("Decompressing %s to %s\n", *inputFile, *outputFile)
		// TODO: Call decompression function
	}
}
