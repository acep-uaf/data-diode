package utility

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "13337"
	CONN_TYPE = "tcp"
)

func TCPServer() {
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)

	if err != nil {
		fmt.Println(">> Error listening: ", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println(">> Server listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(">> Error accepting: ", err.Error())
			return
		}

		go RequestHandler(conn)
	}
}

func Client(ip string, port int) {
	// Create a socket

	client, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), time.Second)

	if err != nil {
		fmt.Println(">> Error establishing connection to the diode input side: ", err.Error())
		log.Fatal(err)
	}
	defer client.Close()

	numberOfSends := 1

	for {
		sendMessage := fmt.Sprintf("This is TCP passthrough test message number: %d", numberOfSends)
		_, err := client.Write([]byte(sendMessage))
		if err != nil {
			fmt.Println(">> Error sending message to the diode input side: ", err.Error())
			log.Fatal(err)
			break
		}

		// if string(response) == "OK\r\n" {
		// 	fmt.Println(">> Message sent successfully!")
		// }

		numberOfSends++

		time.Sleep(1 * time.Second)
	}
}

func Server(ip string, port int) {
	// Begin listening for incoming connections

	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))

	if err != nil {
		fmt.Println(">> Error listening for incoming connections: ", err.Error())
		return
	}
	defer server.Close()

	fmt.Printf(">> Server listening on %s:%d\n", ip, port)

	for {
		// Wait for connection
		connection, err := server.Accept()

		if err != nil {
			fmt.Println(">> Error accepting connection: ", err.Error())
			return
		}

		fmt.Println("Connected to client IP:", connection.RemoteAddr().String())

		go CommunicationHandler(connection)

	}

}

func CommunicationHandler(connection net.Conn) {

	defer connection.Close()

	// Buffer for incoming data (holding recieved data)
	buffer := make([]byte, 10240)

	for {
		// Read incoming data into buffer
		bytesRead, err := connection.Read(buffer)
		if err != nil {
			fmt.Println(">> Error reading: ", err.Error())
			break
		}

		if bytesRead > 0 {
			fmt.Println(">> Message recieved: ", string(buffer[:bytesRead]))
		}

		if bytesRead < 10240 {
			break
		}
	}

}

func RequestHandler(conn net.Conn) {
	buffer := make([]byte, 1024)

	_, err := conn.Read(buffer)

	if err != nil {
		fmt.Println(">> Error reading: ", err.Error())
	}

	conn.Write([]byte("Message received."))

	conn.Close()
}
