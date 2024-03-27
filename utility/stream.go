package utility

import (
	"crypto/md5"
	"bufio"
	"encoding/json"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type Message struct {
	Index     int
	Topic     string
	Payload   string
	Checksum  string
	UUID      string
	Timestamp time.Time
}

var (
	counterMutex   sync.Mutex
	messageCounter int
)

func Craft(topic, payload string) Message {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	uuid := uuid.New().String()

	// TODO: Independent of the topic, the message counter should be incremented?
	messageCounter++

	return Message{
		Index:     messageCounter,
		Topic:     topic,
		Payload:   payload,
		Checksum:  Verification(payload),
		UUID:      uuid,
		Timestamp: time.Now(),
	}
}

func Observability(server string, port int, topic string, message string) error {
	broker := fmt.Sprintf("tcp://%s:%d", server, port)
	clientID := "go_mqtt_client"

	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf(">> Failed to connect to the broker: %v", token.Error())
	}

	defer client.Disconnect(250) // ms

	sample := Craft(topic, message)

	jsonPackage, err := json.Marshal(sample)
	if err != nil {
		return fmt.Errorf(">> Failed to marshal the message: %v", err)
	}

	token := client.Publish(topic, 0, false, jsonPackage)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf(">> Failed to publish the message: %v", token.Error())
	}

	return nil
}

func Republisher(server string, port int, topic string, message string) {
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

func Subscription(server string, port int, topic string, host string, destination int) {
	fmt.Println(">> Example Broker Activity")
	fmt.Println(">> Broker: ", server)
	fmt.Println(">> Port: ", port)

	// MQTT Broker / Client
	url := fmt.Sprintf("tcp://%s:%d", server, port)
	opts := mqtt.NewClientOptions().AddBroker(url)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(">> [!] Failed to connect to the broker: ", token.Error())
	}

	// Callback Function (Incoming Messages)
	handleMessage := func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf(">> Received message on topic: '%s': %s\n", msg.Topic(), msg.Payload())

		// Connection Establishment (Target Host)
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, destination))
		if err != nil {
			fmt.Println(">> [!] Error connecting to the target host: ", err)
			return
		}
		defer conn.Close()

		// Data Transmission
		_, err = conn.Write(msg.Payload())
		if err != nil {
			fmt.Println(">> [!] Error writing to the target host: ", err)
			return
		}
	}

	// Subscription (Topic)
	if token := client.Subscribe(topic, 0, handleMessage); token.Wait() && token.Error() != nil {
		if token.Error() != nil {
			fmt.Println(">> [!] Error subscribing to the topic: ", token.Error())
		}
	}

	// Client Shutdown (SIGINT)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	client.Unsubscribe(topic)
	client.Disconnect(250) // ms
}

func EncapsulatePayload(message string) {
	example := base64.StdEncoding.EncodeToString([]byte(message))
	fmt.Println(">> Encoded message: ", example)
}

func DetectComplete(message string) {
	fmt.Println(">> Detecting complete message...")
}

func ReceivePayload() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Println(">> Received: ", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(">> [!] Error reading from stdin: ", err)
	}
}

func UnencapsulatePayload(message string) {
	decoded, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		fmt.Println(">> [!] Error decoding the message: ", err)
		return
	}
	fmt.Println(">> Decoded message: ", string(decoded))
}

func PublishPayload(message string) {
	fmt.Println(">> Publishing to MQTT broker...")
}

func Verification(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
