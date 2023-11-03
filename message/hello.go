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


	// Request Message
	message, err := example.Hello("Solomon")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(message)
}
