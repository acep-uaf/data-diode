package main

import (
	"testing"
)

func TestCLI(t *testing.T) {
	got := "diode"
	want := "diode"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConfiguration(t *testing.T) {
	got := Configuration{}
	want := Configuration{}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
