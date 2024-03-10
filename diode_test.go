package main

import (
	"fmt"
	"os"
	"testing"

	insights "github.com/acep-uaf/data-diode/insights"
)

func TestCLI(t *testing.T) {
	got := "diode"
	want := "diode"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConfiguration(t *testing.T) {
	_, err := os.Stat("config/settings.yaml")
	if os.IsNotExist(err) {
		t.Errorf("[!] config.yaml does not exist")
	}
}

func TestFileContents(t *testing.T) {
	got := fmt.Sprintf("%x", insights.Checksum())
	want := "477076c6fd8cf48ff2d0159b22bada27588c6fa84918d1c4fc20cd9ddd291dbd"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRepublishContents(t *testing.T) {
	// TODO: Mock the MQTT client.
}
