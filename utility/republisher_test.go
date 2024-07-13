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

		// ? Unique timestamp comparison values with appropriate precision.

		delete(actual, "time")
		delete(expected, "time")

		actualStr, _ := json.Marshal(actual)
		expectedStr, _ := json.Marshal(expected)

		if string(actualStr) != string(expectedStr) {
			t.Errorf("Expected %s but got %s", expectedStr, actualStr)
		}
	}
}
