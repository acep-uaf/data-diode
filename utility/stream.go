package utility

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Message struct {
	Topic    string
	Length   int
	Payload  string
	Checksum string
}

const (
	start   = "###START"
	end     = "###END"
	delimit = false
)

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	DetectContents(string(msg.Payload()), msg.Topic())
}

func Republisher(server string, port int, topic string) {
	location := fmt.Sprintf("tcp://%s:%d", server, port)
	opts := mqtt.NewClientOptions().AddBroker(location).SetClientID("diode_republisher")
	opts.SetDefaultPublishHandler(messageHandler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(">> [!] Failed to connect to the broker: ", token.Error())
	}

	// Callback Function (Incoming Messages)
	handleMessage := func(client mqtt.Client, msg mqtt.Message) {
		DetectContents(string(msg.Payload()), msg.Topic())
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

func EncapsulateContents() {
	scanner := bufio.NewScanner(os.Stdin)
	var line strings.Builder

	// ! Standard Input (Multiple Lines)

	for scanner.Scan() {
		payload := scanner.Text()
		line.WriteString(payload)
	}

	DetectContents(line.String(), "stdin")
}

func DetectContents(message string, topic string) {
	initialize, terminate := start, end

	complete := Message{
		Topic:    topic,
		Length:   len(message),
		Payload:  EncapsulatePayload(message),
		Checksum: Verification(message),
	}

	jsonPackage, err := json.Marshal(complete)
	if err != nil {
		fmt.Println(">> [!] Error marshalling the message: ", err)
		return
	}

	if delimit == true {
		fmt.Println(initialize + string(jsonPackage) + terminate)
	} else {
		fmt.Println(string(jsonPackage))
	}
}

func EncapsulatePayload(message string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(message))
	return encoded
}

func UnencapsulatePayload(message string) {
	// TODO: Test case(s) for various message lengths and content.

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
