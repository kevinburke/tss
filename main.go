// Command tss prints timestamps relative to the program start, and the previous
// line of input.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tss "github.com/kevinburke/tss/lib"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `tss [-v] [-h]

Annotate stdin with timestamps per line.
`)
	}
}

const Version = "0.3"

func main() {
	version := flag.Bool("version", false, "Print the version string")
	v := flag.Bool("v", false, "Print the version string")
	flag.Parse()
	if *version || *v {
		fmt.Fprintf(os.Stderr, "tss version %s\n", Version)
		os.Exit(2)
	}
	if _, err := tss.Copy(os.Stdout, os.Stdin); err != nil {
		log.Fatal(err)
	}
}
