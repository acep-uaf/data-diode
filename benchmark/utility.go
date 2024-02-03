package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/guptarohit/asciigraph"
)

func checksum() {
	// https://en.wikipedia.org/wiki/Rolling_hash

	passthrough := []string{
		"hello",
		"world",
	}

	fmt.Println(">> Words: ", passthrough)

	for _, word := range passthrough {
		sum := sha256.Sum256([]byte(word))
		fmt.Printf("%x\n", sum)
	}
}

func example() {
	definition := []float64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144}
	graph := asciigraph.Plot(definition)

	fmt.Println(graph)
}

func main() {
	checksum()
	example()
}
