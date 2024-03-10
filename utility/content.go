package utility

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type LineContent struct {
	Number  int
	Content string
}

type FileContent struct {
	Lines map[int]string
}

type Readability struct {
	Words      int
	Characters int
	Paragraphs int
	Sentences  int
}

func ReadLineContent(location string) FileContent {
	file, err := os.Open(location)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := make(map[int]string)

	scanner := bufio.NewScanner(file)

	lineNumber := 1

	for scanner.Scan() {
		lineContent := scanner.Text()
		lines[lineNumber] = lineContent
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return FileContent{Lines: lines}
}

func OutputStatistics(content FileContent) {
	// ? Contextual information about the file content.
	fmt.Println(">> Number of lines: ", len(content.Lines))
}

func PrintFileContent(content FileContent) {
	for i := 1; i <= len(content.Lines); i++ {
		fmt.Println(">> ", content.Lines[i])
	}
}
