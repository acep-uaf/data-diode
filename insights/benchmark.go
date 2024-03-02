package analysis

import (
	"crypto/sha256"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func Checksum() [32]byte {
	location := "https://www.gutenberg.org/cache/epub/2701/pg2701.txt"

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

func Value() int {
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	return randGenerator.Intn(100)
}
