package utility

import (
	"fmt"
	"net"
	"time"
)

const (
	ACKNOWLEDGEMENT = "OK\r\n"
	CONN_TYPE       = "tcp"
	MAX_ATTEMPTS    = 2
	CHUNK_SIZE      = 1460  // ? Characters
	SAMPLE          = 10240 // 10 Kbytes
)

func StartPlaceholderClient(host string, port int) {

	for i := 1; i <= MAX_ATTEMPTS; i++ {
		fmt.Printf(">> [%d of %d] Dialing host %s on port %d via %s...\n", i, MAX_ATTEMPTS, host, port, CONN_TYPE)
	}

	conn, err := net.Dial(CONN_TYPE, fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		fmt.Println(">> [!] Error connecting: ", err.Error())
		return
	}
	defer conn.Close()

	message := "The quick brown fox jumps over the lazy dog.\n"

	for try := 1; try <= MAX_ATTEMPTS; try++ {
		if len(message) > CHUNK_SIZE {
			index := 0

			for index < len(message) {
				chunk := message[index : index+CHUNK_SIZE]

				// I. Send Chunk

				_, err := conn.Write([]byte(chunk))
				if err != nil {
					fmt.Println(">> [!] Error sending data: ", err)
					return
				}

				// II. Wait for ACK

				response := make([]byte, 4)
				_, err = conn.Read(response)
				if err != nil {
					fmt.Println(">> [!] Error receiving ACK: ", err)
					return
				}

				// III. Diode Response

				if string(response) != ACKNOWLEDGEMENT {
					fmt.Println(">> [?] Invalid ACK received.")
					return
				}

				fmt.Printf(">> Successfully sent message to diode: %s\n", chunk)

				index += CHUNK_SIZE
			}
		}

		time.Sleep(1 * time.Second)
	}

	buffer := make([]byte, SAMPLE)

	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(">> [!] Error reading response: ", err.Error())
		return
	}

	fmt.Printf(">> Server response: %s\n", string(buffer[:bytesRead]))
}

func StartPlaceholderServer(host string, port int) {
	listener, err := net.Listen(CONN_TYPE, fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println(">> [!] Error listening: ", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println(">> Server listening on: ", listener.Addr())

	for {
		fmt.Println(">> Server waiting for connection...")

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(">> [!] Error accepting: ", err.Error())
			return
		}

		fmt.Println(">> Server accepted connection from: ", conn.RemoteAddr())

		go func(conn net.Conn) {
			defer conn.Close()

			for {
				data := make([]byte, SAMPLE)
				bytesRead, err := conn.Read(data)
				if err != nil {
					fmt.Println(">> [!] Error receiving data: ", err.Error())
					if err.Error() == "EOF" {
						fmt.Println(">> Connection closed by client.")
						return
					}
				}

				fmt.Printf(">> Received data: %s\n", string(data[:bytesRead]))

				_, err = conn.Write([]byte(ACKNOWLEDGEMENT))
				if err != nil {
					fmt.Println(">> [!] Error sending ACK: ", err.Error())
					return
				}
			}
		}(conn)
	}
}
