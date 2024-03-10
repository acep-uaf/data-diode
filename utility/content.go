package utility

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type FileContent struct {
	Content []string
}

func ReadFileContent(location string) FileContent {
	file, err := os.Open(location)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var content []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return FileContent{Content: content}
}

func PrintFileContent(content FileContent) {
	for _, line := range content.Content {
		fmt.Println(">> ", line)
	}
}
