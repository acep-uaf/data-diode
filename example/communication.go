package example

import (
	"errors"
	"fmt"
	"math/rand"
)

// Communication returns a greeting for the named person.
func Communication(name string) (string, error) {
	// If no name was given, return an error with a message.
	if name == "" {
		return "", errors.New("empty name")
	}
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf(randomFormat(), name)
	return message, nil
}

func Communications(names []string) (map[string]string, error) {
    messages := make(map[string]string)
    for _, name := range names {
        message, err := Communication(name)
        if err != nil {
            return nil, err
        }
        messages[name] = message
    }
    return messages, nil
}

func randomFormat() string {
	formats := []string{
		"First sample message: %v.",
		"Second sample message: %v.",
		"Third sample message: %v.",
	}

	return formats[rand.Intn(len(formats))]
}
