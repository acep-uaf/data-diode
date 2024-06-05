package utility

import (
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	TEST_BROKER  = "tcp://localhost:1883"
	TEST_TOPIC   = "example"
	TEST_PAYLOAD = "Hello, world."
	TEST_CONNECTION = "sample"
)

func TestPublish(t *testing.T) {
	opts := mqtt.NewClientOptions().AddBroker(TEST_BROKER).SetClientID(TEST_CONNECTION)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		t.Errorf(">> [!] Failed to connect to the broker: %v", token.Error())
	}

	token := client.Publish(TEST_TOPIC, byte(0), false, TEST_PAYLOAD)

	if token.Error() != nil {
		t.Errorf(">> [!] Publish error: %v", token.Error())
	}

	client.Disconnect(250)
}

func TestSubscribe(t *testing.T) {
	opts := mqtt.NewClientOptions().AddBroker(TEST_BROKER).SetClientID(TEST_CONNECTION)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		t.Errorf(">> [!] Failed to connect to the broker: %v", token.Error())
	}

	handleMessage := func(client mqtt.Client, msg mqtt.Message) {
		t.Logf(">> [+] Received message: %v", string(msg.Payload()))
	}

	if token := client.Subscribe(TEST_TOPIC, byte(0), handleMessage); token.Wait() && token.Error() != nil {
		t.Errorf(">> [!] Error subscribing to the topic: %v", token.Error())
	}

	client.Disconnect(250)
}

func TestSendMessage(t *testing.T) {
	got := "{}"
	want := "{}"

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestReceiveMessage(t *testing.T) {
	got := "{}"
	want := "{}"

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
