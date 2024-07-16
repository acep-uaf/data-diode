package utility

import (
	"fmt"
	"testing"
)

func TestMakeTimestamp(t *testing.T) {
	actual := MakeTimestamp()
	expected := 16

	length := len(fmt.Sprintf("%d", actual))

	if length != expected {
		t.Errorf("Expected %d but got %d", expected, length)
	}
}

func TestVerification(t *testing.T) {
	testcases := []struct {
		input string
	}{
		{"data-diode"},
	}

	for _, test := range testcases {
		actual := Verification(test.input)
		expected := "ef217cdf54b0b3ece89a0d7686da550b"

		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	}
}

func TestEncapsulatePayload(t *testing.T) {
	testcases := []struct {
		input string
	}{
		{"data-diode"},
	}

	for _, test := range testcases {
		actual := EncapsulatePayload(test.input)
		expected := "ZGF0YS1kaW9kZQ=="

		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	}
}

func TestUnencapsulatePayload(t *testing.T) {
	testcases := []struct {
		input string
	}{
		{"ZGF0YS1kaW9kZQ=="},
	}

	for _, test := range testcases {
		actual := UnencapsulatePayload(test.input)
		expected := "data-diode"

		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	}
}
