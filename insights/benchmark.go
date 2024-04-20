package analysis

import (
	"crypto/sha256"
	"io"
	"net/http"
)

func Pong() bool {
	location := "https://example.com/"

	resp, err := http.Get(location)

	if err != nil {
		return false
	}
	defer resp.Body.Close()

	println(">> Status: ", resp.Status)

	return true
}

func Saturate() {
	println(">> Data Diode Components")
}

func Checksum() [32]byte {
	location := "https://www.gutenberg.org/cache/epub/84/pg84.txt"

	resp, err := http.Get(location)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// https://en.wikipedia.org/wiki/Rolling_hash

	hash := sha256.Sum256([]byte(body))

	return hash
}
