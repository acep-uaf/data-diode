package utility

import (
	"fmt"
	"math/rand"
	"net"
)

const (
	CONN_TYPE = "tcp"
	SAMPLE    = 1024
)

func StartPlaceholderClient(CONN_HOST string, CONN_PORT string, CONN_TYPE string) {
	upperBound := rand.Intn(1) + 1

	for i := 1; i <= upperBound; i++ {
		fmt.Printf(">> [%d of %d] Dialing host %s on port %s via %s...\n", i, upperBound, CONN_HOST, CONN_PORT, CONN_TYPE)
	}

	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT) // Server Connection

	if err != nil {
		fmt.Println(">> [!] Error connecting: ", err.Error())
		return
	}
	defer conn.Close()

	message := "The quick brown fox jumps over the lazy dog.\n"

	_, err = conn.Write([]byte(message))

	buffer := make([]byte, SAMPLE)

	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(">> [!] Error reading response: ", err.Error())
		return
	}

	fmt.Printf(">> Server response: %s\n", string(buffer[:bytesRead]))
}

func StartPlaceholderServer(CONN_HOST string, CONN_PORT string, CONN_TYPE string) {
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)

	if err != nil {
		fmt.Println(">> [!] Error listening: ", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println(">> Server listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(">> [!] Error accepting: ", err.Error())
			return
		}

		go RequestHandler(conn)
	}
}

func RequestHandler(conn net.Conn) {
	buffer := make([]byte, SAMPLE)

	_, err := conn.Read(buffer)

	if err != nil {
		fmt.Println(">> [!] Error reading: ", err.Error())
	}

	conn.Write([]byte("[✅] Message received.\n"))

	fmt.Printf(">> Message received: %s\n", string(buffer))

	conn.Close()
}
