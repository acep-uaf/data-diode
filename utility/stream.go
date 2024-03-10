package utility

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Message struct {
	Index    int
	Topic    string
	Payload  string
	Checksum string
}

var (
	counterMutex   sync.Mutex
	messageCounter int
)

func Craft(topic, payload string) Message {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	// TODO: Independent of the topic, the message counter should be incremented?
	messageCounter++

	return Message{
		Index:    messageCounter,
		Topic:    topic,
		Payload:  payload,
		Checksum: Verification(payload),
	}
}

func Observability(server string, port int, topic string, message string) {
	broker := fmt.Sprintf("tcp://%s:%d", server, port)
	clientID := "go_mqtt_client"

	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	defer client.Disconnect(250) // ms

	sample := Craft(topic, message)

	jsonPackage, err := json.Marshal(sample)
	if err != nil {
		panic(err)
	}

	token := client.Publish(topic, 0, false, jsonPackage)
	token.Wait()
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

func Verification(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
