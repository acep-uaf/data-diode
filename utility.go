package main

import (
	"fmt"
)

func example() {
	// https://en.wikipedia.org/wiki/Rolling_hash
	passthrough := []string{
		"hello",
		"world",
	}

	fmt.Println(">> Words:", passthrough)
}
