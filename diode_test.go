package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	insights "github.com/acep-uaf/data-diode/insights"
	"github.com/acep-uaf/data-diode/utility"
)

var (
	BackupConfiguration  = "config/B4-0144-355112.json"
	SystemSettings       = "config/settings.yaml"
	ProjectDocumentation = "docs/SOP.md"
	FileChecksum         = "477076c6fd8cf48ff2d0159b22bada27588c6fa84918d1c4fc20cd9ddd291dbd"
	SampleMessage        = "Hello, world."
	InterfaceSize        = 1024
	InterfaceProtocol    = "tcp"
	InterfaceAddress     = "localhost:13337"
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
	binary := exec.Command("go", "build", "-o", "diode")
	buildErr := binary.Run()
	if buildErr != nil {
		t.Fatalf("[!] Failed to build CLI binary: %v", buildErr)
	}

	cmd := exec.Command("./diode")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("[!] Failed to execute CLI command: %v", err)
	}

	expectation := "diode: try 'diode --help' for more information"
	reality := strings.TrimSpace(stdout.String())
	if reality != expectation {
		t.Errorf("[?] Expected output: %q, but got: %q", expectation, reality)
	}

	if stderr.Len() > 0 {
		t.Errorf("[?] Unexpected error output: %q", stderr.String())
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
		t.Errorf(">> got %q, want %q", got, want)
	}
}

func TestBinaryContents(t *testing.T) {
	// TODO: Implement the following:
	// - Craft a text message containing binary data + checksum.
	// - Ensure transmission across data diode without corrupted information.
	// - Check for uuenconding and base64 encoding / delimiters.

	sample := []byte(SampleMessage)

	if len(sample) == 0 {
		t.Errorf("[!] No binary contents...")
	}
}

func TestEchoMessage(t *testing.T) {
	go func() {
		listener, err := net.Listen(InterfaceProtocol, InterfaceAddress)
		if err != nil {
			t.Errorf("[!] Failed to start TCP server: %v", err)
		}
		defer listener.Close()

		conn, err := listener.Accept()
		if err != nil {
			t.Errorf("[!] Failed to accept connection: %v", err)
		}
		defer conn.Close()

		buf := make([]byte, InterfaceSize)
		n, err := conn.Read(buf)
		if err != nil {
			t.Errorf("[!] Failed to read message: %v", err)
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			t.Errorf("[!] Failed to write message: %v", err)
		}
	}()

	// TODO: Mock the TCP client/server to simulate the transmission of data.

	conn, err := net.Dial(InterfaceProtocol, InterfaceAddress)
	if err != nil {
		t.Fatalf("[!] Failed to connect to TCP server: %v", err)
	}
	defer conn.Close()

	message := SampleMessage
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("[!] Failed to send message: %v", err)
	}

	buf := make([]byte, len(message))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("[!] Failed to read echoed message: %v", err)
	}

	match := string(buf[:n])
	if match != message {
		t.Errorf("[!] Echoed message does not match original message: got %q, want %q", match, message)
	}
}

func TestRepublishContents(t *testing.T) {
	location := ProjectDocumentation
	broker := "localhost"
	topic := "test/message"
	port := 1883

	// TODO: Mock the MQTT connection.

	utility.RepublishContents(location, broker, topic, port)

	if len(location) == 0 {
		t.Errorf("[!] No location specified...")
	}

	if len(broker) == 0 {
		t.Errorf("[!] No broker specified...")
	}

	if len(topic) == 0 {
		t.Errorf("[!] No topic specified...")
	}

	if port == 0 {
		t.Errorf("[!] No port specified...")
	}
}

func TestRepublishContents(t *testing.T) {
	// TODO: Mock the MQTT client.
}
