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

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/urfave/cli/v2"
	"rsc.io/quote"
)

func newClient(diodeInputSideIP string, diodeTcpPassthroughPort int) {
	// Create a socket

	client, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", diodeInputSideIP, diodeTcpPassthroughPort), time.Second)

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

func newServer(targetTcpServerIP string, targetTcpServerPort int) {
	// Begin listening for incoming connections

	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", targetTcpServerIP, targetTcpServerPort))

	if err != nil {
		fmt.Println(">> Error listening for incoming connections: ", err.Error())
		return
	}
	defer server.Close()

	fmt.Printf(">> Server listening on %s:%d\n", targetTcpServerIP, targetTcpServerPort)

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

func sampleMetrics() {
	fmt.Println(">> Local time: ", time.Now())
	fmt.Println(">> UTC time: ", time.Now().UTC())
}

func demo() {
	mqttBrokerIP := "test.mosquitto.org"
	mqttBrokerPort := 1883
	mqttBrokerMessage := "Hello, world."
	mqttBrokerTopic := "test/message"

	fmt.Println(">> MQTT")
	fmt.Println(">> Broker: ", mqttBrokerIP)
	fmt.Println(">> Port: ", mqttBrokerPort)

	// Source: https://github.com/eclipse/paho.mqtt.golang/blob/master/cmd/simple/main.go
	var example mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf(">> Topic: %s\n", msg.Topic())
		fmt.Printf(">> Message: %s\n", msg.Payload())
	}

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	// Initial Connection
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBrokerIP, mqttBrokerPort))
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(example)
	opts.SetPingTimeout(1 * time.Second)

	// Create and start a client using the above ClientOptions
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to a topic
	if token := client.Subscribe(mqttBrokerTopic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Publish to a topic
	token := client.Publish(mqttBrokerTopic, 0, false, mqttBrokerMessage)
	token.Wait()

	time.Sleep(6 * time.Second)

	// Disconnect from the broker
	if token := client.Unsubscribe(mqttBrokerTopic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	client.Disconnect(250)

	time.Sleep(1 * time.Second)

}

func main() {

	// Configuration Options

	diodeInputSideIP := "192.168.1.99"
	diodeTcpPassthroughPort := 50000

	targetTcpServerIP := "192.168.1.20"
	targetTcpServerPort := 503

	app := &cli.App{
		Name:  "diode",
		Usage: "A command line tool for interacting with data diodes.",
		Action: func(cCtx *cli.Context) error {
			fmt.Println(quote.Go())
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "Input side of the data diode",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("----- INPUT -----")
					newClient(diodeInputSideIP, diodeTcpPassthroughPort)
					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Output side of the data diode",
				Action: func(sCtx *cli.Context) error {
					fmt.Println("----- OUTPUT -----")
					newServer(targetTcpServerIP, targetTcpServerPort)
					return nil
				},
			},
			{
				Name:    "diagnostics",
				Aliases: []string{"d"},
				Usage:   "Debug diagnostics via configuration settings",
				Action: func(dCtx *cli.Context) error {
					fmt.Println("----- DIAGNOSTICS -----")
					input := fmt.Sprintf("%s:%d", diodeInputSideIP, diodeTcpPassthroughPort)
					output := fmt.Sprintf("%s:%d", targetTcpServerIP, targetTcpServerPort)
					fmt.Println(">> Client: ", input)
					fmt.Println(">> Server: ", output)
					return nil
				},
			},
			{
				Name:    "benchmark",
				Aliases: []string{"b"},
				Usage:   "System benchmark analysis + report performance metrics",
				Action: func(bCtx *cli.Context) error {
					fmt.Println("----- BENCHMARKS -----")
					sampleMetrics()
					example()
					return nil
				},
			},
			{
				Name:    "mqtt",
				Aliases: []string{"m"},
				Usage:   "MQTT (republisher) demo",
				Action: func(mCtx *cli.Context) error {
					fmt.Println("----- MQTT -----")
					demo()
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the version of the diode CLI",
				Action: func(vCtx *cli.Context) error {
					fmt.Println(">> diode version 0.0.4")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal()
	}
}
