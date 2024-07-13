package utility

import (
	"net"
	"testing"
)

func TestSendMessage(t *testing.T) {
	server, err := net.Listen(CONN_TYPE, "localhost:0")
	if err != nil {
		t.Fatalf(">> [!] Failed to create mock server: %v", err)
	}
	defer server.Close()

	go func() {
		conn, err := server.Accept()
		if err != nil {
			t.Error(">> [!] Failed to accept connection:", err)
			return
		}
		defer conn.Close()

		buffer := make([]byte, CHUNK_SIZE)
		_, err = conn.Read(buffer)
		if err != nil {
			t.Error("[!] Failed to read data:", err)
			return
		}

		_, err = conn.Write([]byte(ACKNOWLEDGEMENT))
		if err != nil {
			t.Error("[!] Failed to send acknowledgement:", err)
			return
		}
	}()

	input := "data-diode"
	client := server.Addr().String()
	SendMessage(input, client)
}
