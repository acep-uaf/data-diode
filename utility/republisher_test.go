package utility

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDetectContents(t *testing.T) {
	testcases := []struct {
		input string
		topic string
	}{
		{"data-diode", "first"},
		{"99775", "second"},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum", "third"},
	}

	for _, test := range testcases {
		actualJSON := DetectContents(test.input, test.topic)
		expectedJSON := fmt.Sprintf(`{"time": %d, "topic": "%s", "b64payload": "%s"}`, MakeTimestamp(), test.topic, EncapsulatePayload(test.input))

		var actual, expected map[string]interface{}
		if err := json.Unmarshal([]byte(actualJSON), &actual); err != nil {
			t.Errorf(">> [!] Failed to unmarshal the actual output: %v", err)
		}
		if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
			t.Errorf(">> [!] Failed to unmarshal the expected output: %v", err)
		}

		// TODO: Unique timestamp comparison values with appropriate precision.

		delete(actual, "time")
		delete(expected, "time")

		actualStr, _ := json.Marshal(actual)
		expectedStr, _ := json.Marshal(expected)

		if string(actualStr) != string(expectedStr) {
			t.Errorf("Expected %s but got %s", expectedStr, actualStr)
		}
	}
}

func TestRepackageContents(t *testing.T) {
	testcases := []struct {
		input  string
		topic  string
		output string
	}{
		{"{\"time\": 1234567890123456, \"topic\": \"test-topic\", \"b64payload\": \"ZGF0YS1kaW9kZQ==\"}", "test-topic", "data-diode"},
	}

	for _, test := range testcases {
		actual := RepackageContents(test.input, test.topic)
		expected := test.output

		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	}
}

func TestPublishPayload(t *testing.T) {
	testcases := []struct {
		server  string
		port    int
		topic   string
		message string
	}{
		{"localhost", 1883, "test-topic", "data-diode"},
	}

	for _, test := range testcases {
		PublishPayload(test.server, test.port, test.topic, test.message)
	}
}
