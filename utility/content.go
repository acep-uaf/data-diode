package utility

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

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

func ExampleContents(location string) {
	sample := ReadLineContent(location)
	PrintFileContent(sample)
	OutputStatistics(sample)
}

func RepublishContents(location string, mqttBrokerIP string, mqttBrokerTopic string, mqttBrokerPort int) error {
	if _, err := os.Stat(location); os.IsNotExist(err) {
		fmt.Println(">> File not found: ", location)
		return err
	}

	fileContent := ReadLineContent(location)

	fmt.Println(">> Server: ", mqttBrokerIP)
	fmt.Println(">> Topic: ", mqttBrokerTopic)
	fmt.Println(">> Port: ", mqttBrokerPort)

	start := time.Now()

	for i := 1; i <= len(fileContent.Lines); i++ {
		Observability(mqttBrokerIP, mqttBrokerPort, mqttBrokerTopic, fileContent.Lines[i])
	}

	t := time.Now()

	elapsed := t.Sub(start)

	if len(fileContent.Lines) == 0 {
		fmt.Println(">> No message content sent.")
	} else if len(fileContent.Lines) == 1 {
		fmt.Println(">> Sent message from ", location, " to topic: ", mqttBrokerTopic, " in ", elapsed)
	} else {
		fmt.Println(">> Sent ", len(fileContent.Lines), " messages from ", location, " to topic: ", mqttBrokerTopic, " in ", elapsed)
	}

	return nil
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
