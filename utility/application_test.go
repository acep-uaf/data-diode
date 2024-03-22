package utility

import "testing"

func TestClient(t *testing.T) {
	got := "server"
	want := "client"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestServer(t *testing.T) {
	got := "client"
	want := "server"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
