package utility

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type InputDiodeMessage struct {
	Time       int    `json:"time"`
	Topic      string `json:"topic"`
	B64Payload string `json:"b64payload"`
}

type OutputDiodeMessage struct {
	Time       int    `json:"time"`
	Topic      string `json:"topic"`
	B64Payload string `json:"b64payload"`
	Payload    string `json:"payload"`
	Length     int    `json:"length"`
	Checksum   string `json:"checksum"`
}

func InboundMessageFlow(server string, port int, topic string, arrival string) {
	location := fmt.Sprintf("tcp://%s:%d", server, port)
	opts := mqtt.NewClientOptions().AddBroker(location).SetClientID("in_rec_msg")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(">> [!] Failed to connect to the broker: ", token.Error())
	}

	handleMessage := func(client mqtt.Client, msg mqtt.Message) {
		contents := DetectContents(string(msg.Payload()), msg.Topic())
		SendMessage(contents, arrival)
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
}

func OutboundMessageFlow(server string, port int, topic string, destination string) {
	example, err := RecieveMessage(destination)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Detect, decode, and unencapsulate the message before publishing.

	specificity := "diode/telemetry"
	PublishPayload(server, port, specificity, example)
}

func DetectContents(message string, topic string) string {
	complete := InputDiodeMessage{
		Time:       int(MakeTimestamp()),
		Topic:      topic,
		B64Payload: EncapsulatePayload(message),
	}

	jsonPackage, err := json.Marshal(complete)
	if err != nil {
		log.Fatalf(">> [!] Error marshalling the message: %v", err)
	}

	return string(jsonPackage)
}

func EncapsulatePayload(message string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(message))
	return encoded
}

func UnencapsulatePayload(message string) string {
	// TODO: Test case(s) for various message lengths and content.

	decoded, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		fmt.Println(">> [!] Error decoding the message: ", err)
	}
	return string(decoded)
}

func PublishPayload(server string, port int, topic string, message string) {
	location := fmt.Sprintf("tcp://%s:%d", server, port)
	opts := mqtt.NewClientOptions().AddBroker(location).SetClientID("diode_republisher")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(">> [!] Failed to connect to the broker: ", token.Error())
	}

	if token := client.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		fmt.Println(">> [!] Error publishing the message: ", token.Error())
	}

	client.Disconnect(250)
}

func MakeTimestamp() int64 {
	return time.Now().UnixMilli()
}

func Verification(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
