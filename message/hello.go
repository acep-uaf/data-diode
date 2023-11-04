package main

import (
	"fmt"
	"log"

	"example.com/example"
)

func main() {
	// Predefined Properties
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	// Name Slice
	names := []string{"A", "B", "C"}

	// Request Message(s)
	messages, err := example.Communications(names)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messages)
}
