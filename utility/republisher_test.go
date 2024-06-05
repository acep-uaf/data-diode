package utility

import (
	"encoding/base64"
	"testing"
	"time"
)

func TestPayloadStructure(t *testing.T) {

	sample := "Hello, world."
	location := "test"

	message := InputDiodeMessage{
		Time:       int(time.Now().Unix()),
		Topic:      location,
		B64Payload: base64.StdEncoding.EncodeToString([]byte(sample)),
	}

	got := message.B64Payload
	want := EncapsulatePayload(sample)

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
