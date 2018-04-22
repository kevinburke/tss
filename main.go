// Command tss prints timestamps relative to the program start, and the previous
// line of input.
package main

import (
	"flag"
	"log"
	"os"

	tss "github.com/kevinburke/tss/lib"
)

func main() {
	flag.Parse()
	if _, err := tss.Copy(os.Stdout, os.Stdin); err != nil {
		log.Fatal(err)
	}
}
