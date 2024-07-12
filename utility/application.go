package utility

import (
	"fmt"
	"log"
	"net"
)

const (
	ACKNOWLEDGEMENT = "OK\r\n"
	CONN_TYPE       = "tcp"
	CHUNK_SIZE      = 1460
	SAMPLE          = 10240
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SendMessage(input string, client string) {
	conn, err := net.Dial(CONN_TYPE, client)
	if err != nil {
		log.Fatalf(">> [!] Error connecting to diode client: %v", err)
	}
	defer conn.Close()

	for index := 0; index < len(input); index += CHUNK_SIZE {
		chunk := input[index:min(index+CHUNK_SIZE, len(input))]

		_, err := conn.Write([]byte(chunk))
		if err != nil {
			log.Fatalf(">> [!] Error sending data: %v", err)
		}

		response := make([]byte, len(ACKNOWLEDGEMENT))
		_, err = conn.Read(response)
		if err != nil {
			log.Fatalf(">> [!] Error receiving ACK: %v", err)
		}

		if string(response) != ACKNOWLEDGEMENT {
			log.Fatalf(">> [?] Invalid ACK received.")
		}

		fmt.Println(chunk)
	}
}

func StartPlaceholderClient(host string, port int) {
	conn, err := net.Dial(CONN_TYPE, fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		fmt.Println(">> [!] Error connecting: ", err.Error())
		return
	}
	defer conn.Close()

	message := "The quick brown fox jumps over the lazy dog.\n"

	ProcessMessage(message, conn)

	buffer := make([]byte, SAMPLE)

	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(">> [!] Error reading response: ", err.Error())
		return
	}

	fmt.Printf(">> Server response: %s\n", string(buffer[:bytesRead]))
}

func ProcessMessage(message string, conn net.Conn) {
	if len(message) > CHUNK_SIZE {
		index := 0

		for index < len(message) {
			chunk := message[index : index+CHUNK_SIZE]

			_, err := conn.Write([]byte(chunk))
			if err != nil {
				fmt.Println(">> [!] Error sending data: ", err)
			}

			response := make([]byte, 4)
			_, err = conn.Read(response)
			if err != nil {
				fmt.Println(">> [!] Error receiving ACK: ", err)
			}

			if string(response) != ACKNOWLEDGEMENT {
				fmt.Println(">> [?] Invalid ACK received.")
			}

			fmt.Printf(">> Successfully sent message to diode: %s\n", chunk)

			index += CHUNK_SIZE
		}
	}
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
					return
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

func RecieveMessage(destination string, messages chan<- string) error {
	server, err := net.Listen("tcp", destination)
	if err != nil {
		fmt.Println(">> [!] Error connecting to diode: ", err)
		return err
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(">> [!] Error accepting connection: ", err)
			continue
		}

		go func(conn net.Conn) {
			message, err := connectionHandler(conn)
			if err != nil {
				fmt.Println(">> [!] Error handling connection: ", err)
				return
			}

			messages <- message
		}(conn)
	}
}

func connectionHandler(conn net.Conn) (string, error) {
	defer conn.Close()

	buffer := make([]byte, SAMPLE)

	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(">> [!] Error reading data: ", err)
		return "", err
	}

	return string(buffer[:bytesRead]), nil
}
