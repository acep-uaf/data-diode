// ______      _         ______ _           _
// |  _  \    | |        |  _  (_)         | |
// | | | |__ _| |_ __ _  | | | |_  ___   __| | ___
// | | | / _` | __/ _` | | | | | |/ _ \ / _` |/ _ \
// | |/ / (_| | || (_| | | |/ /| | (_) | (_| |  __/
// |___/ \__,_|\__\__,_| |___/ |_|\___/ \__,_|\___|

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/acep-uaf/data-diode/utility"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var (
	SemVer string
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

type Configuration struct {
	Input struct {
		IP   string
		Port int
	}
	Output struct {
		IP   string
		Port int
	}
	Broker struct {
		Server  string
		Port    int
		Topic   string
		Message string
	}
}

func newTCPServer() {
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

		go requestHandler(conn)
	}
}

func newClient(ip string, port int) {
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

func newServer(ip string, port int) {
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

		go communicationHandler(connection)

	}

}

func communicationHandler(connection net.Conn) {

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

func requestHandler(conn net.Conn) {
	buffer := make([]byte, 1024)

	_, err := conn.Read(buffer)

	if err != nil {
		fmt.Println(">> Error reading: ", err.Error())
	}

	conn.Write([]byte("Message received."))

	conn.Close()
}

func sampleMetrics(server string, port int) {
	fmt.Println(">> Local time: ", time.Now())
	fmt.Println(">> UTC time: ", time.Now().UTC())
	fmt.Println(">> Value: ", utility.Value())
}

func demoRepublisher(server string, port int, topic string, message string) {
	fmt.Println(">> MQTT")
	fmt.Println(">> Broker: ", server)
	fmt.Println(">> Port: ", port)

	// Source: https://github.com/eclipse/paho.mqtt.golang/blob/master/cmd/simple/main.go
	var example mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf(">> Topic: %s\n", msg.Topic())
		fmt.Printf(">> Message: %s\n", msg.Payload())
	}

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	// Initial Connection
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", server, port))
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(example)
	opts.SetPingTimeout(1 * time.Second)

	// Create and start a client using the above ClientOptions
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to a topic
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Publish to a topic
	token := client.Publish(topic, 0, false, message)
	token.Wait()

	time.Sleep(6 * time.Second)

	// Disconnect from the broker
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	client.Disconnect(250)

	time.Sleep(1 * time.Second)

}

func main() {
	data, err := os.ReadFile("config.yaml")

	if err != nil {
		panic(err)
	}

	var config Configuration

	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	// Configuration Settings

	diodeInputSideIP := config.Input.IP
	diodeTCPPassthroughPort := config.Input.Port
	targetTCPServerIP := config.Output.IP
	targetTCPServerPort := config.Output.Port

	mqttBrokerIP := config.Broker.Server
	mqttBrokerPort := config.Broker.Port
	mqttBrokerMessage := config.Broker.Message
	mqttBrokerTopic := config.Broker.Topic

	app := &cli.App{
		Name:  "diode",
		Usage: "Tool for interacting with data diode(s) via command-line interface (CLI).",
		Action: func(cCtx *cli.Context) error {
			fmt.Println("diode: try 'diode --help' for more information")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "Input side of the data diode",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("----- INPUT -----")
					newClient(diodeInputSideIP, diodeTCPPassthroughPort)
					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Output side of the data diode",
				Action: func(sCtx *cli.Context) error {
					fmt.Println("----- OUTPUT -----")
					newServer(targetTCPServerIP, targetTCPServerPort)
					return nil
				},
			},
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "Testing state synchronization via diode I/O",
				Action: func(tCtx *cli.Context) error {
					fmt.Println("----- TEST -----")
					newTCPServer()
					return nil
				},
			},
			{
				Name:    "diagnostics",
				Aliases: []string{"d"},
				Usage:   "Debug diagnostics via configuration settings",
				Action: func(dCtx *cli.Context) error {
					fmt.Println("----- DIAGNOSTICS -----")
					fmt.Printf("%+v\n", config)
					return nil
				},
			},
			{
				Name:    "benchmark",
				Aliases: []string{"b"},
				Usage:   "System benchmark analysis + report performance metrics",
				Action: func(bCtx *cli.Context) error {
					fmt.Println("----- BENCHMARKS -----")
					sampleMetrics(CONN_HOST, 3333)
					return nil
				},
			},
			{
				Name:    "mqtt",
				Aliases: []string{"m"},
				Usage:   "MQTT (TCP stream) demo",
				Action: func(mCtx *cli.Context) error {
					fmt.Println("----- MQTT -----")
					demoRepublisher(mqttBrokerIP, mqttBrokerPort, mqttBrokerTopic, mqttBrokerMessage)
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the version of the diode CLI",
				Action: func(vCtx *cli.Context) error {
					fmt.Println(">> diode version " + SemVer)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal()
	}
}
