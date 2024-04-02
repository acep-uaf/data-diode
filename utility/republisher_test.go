package utility

import "testing"

func TestSubscribe(t *testing.T) {
	got := "pub"
	want := "sub"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPublish(t *testing.T) {
	got := "sub"
	want := "pub"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
