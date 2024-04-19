package utility

import (
	"net"
	"testing"
)

func TestRecieveMessage(t *testing.T) {
	serverMock, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("[!] Failed to start the mock server: %v", err)
	}
	defer serverMock.Close()

	serverMockAddress := serverMock.Addr().String()

	contents := make(chan string)
	go func() {
		err := RecieveMessage(serverMockAddress, contents)
		if err != nil {
			t.Errorf("[?] Returned an error: %v", err)
		}
	}()

	conn, err := net.Dial("tcp", serverMockAddress)
	if err != nil {
		t.Fatalf("[!] Failed to connect to the mock server: %v", err)
	}
	defer conn.Close()
}
