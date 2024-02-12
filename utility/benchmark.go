package utility

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

func Checksum() {
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

func Value() int {
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	return randGenerator.Intn(100)
}
