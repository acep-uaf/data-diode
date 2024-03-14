package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	insights "github.com/acep-uaf/data-diode/insights"
)

var (
	BackupConfiguration = "config/B4-0144-355112.json"
	SystemSettings      = "config/settings.yaml"
	FileChecksum        = "477076c6fd8cf48ff2d0159b22bada27588c6fa84918d1c4fc20cd9ddd291dbd"
)

type TCP struct {
	ClientTargetIP      string
	ClientTargetPort    int
	ServerTargetIP      string
	ServerPort          int
	ServerSocketTimeout int
}

func TestAPI(t *testing.T) {
	jsonFile, err := os.Open(BackupConfiguration)

	schema := "CAMIO.2024.1.0"
	version := filepath.Base(jsonFile.Name())

	// FIXME: Cross reference the JSON contents, schema version, & configuration file?
	fmt.Println(version, schema)

	if err != nil {
		t.Errorf("[?] %s via %s", err, jsonFile.Name())
	}
}

func TestCLI(t *testing.T) {
	got := "diode"
	want := "diode"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConfiguration(t *testing.T) {
	_, err := os.Stat(SystemSettings)
	if os.IsNotExist(err) {
		t.Errorf("[!] config.yaml does not exist")
	}
}

func TestFileContents(t *testing.T) {
	got := fmt.Sprintf("%x", insights.Checksum())
	want := FileChecksum

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRepublishContents(t *testing.T) {
	// TODO: Mock the MQTT client.
}
