package main

import (
	"crypto/sha256"
	"fmt"
)

func example() {
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
