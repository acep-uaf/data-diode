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
	messages := make(chan string)
	go func() {
		err := RecieveMessage(destination, messages)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	for message := range messages {
		repackaged := RepackageContents(message, topic)
		PublishPayload(server, port, topic, repackaged)
	}
}

func DetectContents(message string, topic string) string {
	complete := InputDiodeMessage{
		Time:       int(MakeTimestamp()),
		Topic:      topic,
		B64Payload: EncapsulatePayload(message),
	}

	jsonPackage, err := json.Marshal(complete)
	if err != nil {
		log.Fatalf(">> [!] Error marshalling the incoming message: %v", err)
	}

	return string(jsonPackage)
}

func RepackageContents(message string, topic string) string {
	var intermediary OutputDiodeMessage
	err := json.Unmarshal([]byte(message), &intermediary)
	if err != nil {
		log.Fatalf(">> [!] Error unmarshalling the message: %v", err)
	}

	// Diode Metadata
	intermediary.Time = int(MakeTimestamp())
	intermediary.Topic = topic
	intermediary.Payload = UnencapsulatePayload(intermediary.B64Payload)
	intermediary.Length = len(intermediary.Payload)
	intermediary.Checksum = Verification(intermediary.Payload)

	// Process Contents
	jsonIntermediary, err := json.Marshal(intermediary)
	if err != nil {
		log.Fatalf(">> [!] Error marshalling the outgoing message: %v", err)
	}

	fmt.Println(string(jsonIntermediary))

	return string(intermediary.Payload)
}

func EncapsulatePayload(message string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(message))
	return encoded
}

func UnencapsulatePayload(message string) string {
	decoded, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		fmt.Println(">> [!] Error decoding the message: ", err)
	}
	return string(decoded)
}

func PublishPayload(server string, port int, topic string, message string) {
	location := fmt.Sprintf("tcp://%s:%d", server, port)
	opts := mqtt.NewClientOptions().AddBroker(location).SetClientID("out_rec_string")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(">> [!] Failed to connect to the broker: ", token.Error())
	}

	if token := client.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		fmt.Println(">> [!] Error publishing the message: ", token.Error())
	}
}

func MakeTimestamp() int64 {
	return time.Now().UnixMicro()
}

func Verification(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
