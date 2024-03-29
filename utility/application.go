package utility

import (
	"fmt"
	"log"
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

		marker := time.Now().UTC().Format(time.RFC3339Nano)
		fmt.Printf(">> [%s] %s\n", marker, chunk)
	}
}

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

	ProcessMessage(message, conn)

	buffer := make([]byte, SAMPLE)

	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(">> [!] Error reading response: ", err.Error())
		return
	}

	fmt.Printf(">> Server response: %s\n", string(buffer[:bytesRead]))
}

func ProcessMessage(message string, conn net.Conn) bool {
	for try := 1; try <= MAX_ATTEMPTS; try++ {
		if len(message) > CHUNK_SIZE {
			index := 0

			for index < len(message) {
				chunk := message[index : index+CHUNK_SIZE]

				_, err := conn.Write([]byte(chunk))
				if err != nil {
					fmt.Println(">> [!] Error sending data: ", err)
					return true
				}

				response := make([]byte, 4)
				_, err = conn.Read(response)
				if err != nil {
					fmt.Println(">> [!] Error receiving ACK: ", err)
					return true
				}

				if string(response) != ACKNOWLEDGEMENT {
					fmt.Println(">> [?] Invalid ACK received.")
					return true
				}

				fmt.Printf(">> Successfully sent message to diode: %s\n", chunk)

				index += CHUNK_SIZE
			}
		}

		time.Sleep(1 * time.Second)
	}
	return false
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

func RecieveMessage(destination string) (string, error) {
	server, err := net.Listen("tcp", destination)
	if err != nil {
		fmt.Println(">> [!] Error connecting to diode: ", err)
		return "", err
	}
	defer server.Close()

	fmt.Println(">> Server listening on: ", server.Addr())

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(">> [!] Error accepting connection: ", err)
			continue
		}

		fmt.Println(">> Server accepted connection from: ", conn.RemoteAddr())

		message, err := connectionHandler(conn)
		if err != nil {
			fmt.Println(">> [!] Error handling connection: ", err)
			break
		}

		return message, nil
	}

	return "", nil
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
